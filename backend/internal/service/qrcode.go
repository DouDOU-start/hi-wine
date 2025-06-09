package service

import (
	"bytes"
	"context"
	"fmt"

	qrcodev1 "backend/api/qrcode/v1"
	"backend/internal/dao"
	"backend/internal/model/entity"
	"backend/internal/utility/minio"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/skip2/go-qrcode"
)

// QrcodeService 二维码服务接口
type QrcodeService interface {
	// CreateTableQrcode 生成桌号二维码
	CreateTableQrcode(ctx context.Context, req *qrcodev1.CreateTableQrcodeReq) (*qrcodev1.TableQrcode, error)

	// GetTableQrcodeList 获取桌号二维码列表
	GetTableQrcodeList(ctx context.Context, req *qrcodev1.TableQrcodeListReq) ([]qrcodev1.TableQrcode, int, error)

	// DeleteTableQrcode 删除桌号二维码
	DeleteTableQrcode(ctx context.Context, id int64) error
}

// 单例实例
var qrcodeServiceInstance = qrcodeService{}

// Qrcode 获取二维码服务实例
func Qrcode() QrcodeService {
	return &qrcodeServiceInstance
}

// 二维码服务实现
type qrcodeService struct{}

// CreateTableQrcode 生成桌号二维码
func (s *qrcodeService) CreateTableQrcode(ctx context.Context, req *qrcodev1.CreateTableQrcodeReq) (*qrcodev1.TableQrcode, error) {
	// 1. 检查桌号是否已存在
	var existQrcode *entity.TableQrcodes
	err := dao.TableQrcodes.Ctx(ctx).
		Where(dao.TableQrcodes.Columns().TableNumber, req.TableNumber).
		Scan(&existQrcode)
	if err != nil {
		return nil, err
	}
	if existQrcode != nil {
		return nil, gerror.New("该桌号已存在")
	}

	// 2. 生成二维码内容
	// 小程序跳转链接，格式：pages/order/index?table_id=xxx
	content := fmt.Sprintf("pages/order/index?table=%s", req.TableNumber)

	// 3. 生成二维码图片
	qr, err := qrcode.New(content, qrcode.Medium)
	if err != nil {
		return nil, gerror.Wrap(err, "生成二维码失败")
	}

	// 设置二维码尺寸
	qr.DisableBorder = false

	// 生成PNG图片数据
	pngData, err := qr.PNG(256)
	if err != nil {
		return nil, gerror.Wrap(err, "生成二维码图片失败")
	}

	// 4. 上传二维码图片到MinIO
	contentType := "image/png"

	// 上传到MinIO
	minioClient := minio.GetClient()
	url, err := minioClient.UploadFileFromReader(
		ctx,
		bytes.NewReader(pngData),
		int64(len(pngData)),
		contentType,
		"qrcodes",
		".png",
	)
	if err != nil {
		return nil, gerror.Wrap(err, "上传二维码图片失败")
	}

	// 5. 保存二维码信息到数据库
	// 使用map直接插入，不设置status字段
	data := g.Map{
		dao.TableQrcodes.Columns().TableNumber: req.TableNumber,
		dao.TableQrcodes.Columns().QrcodeUrl:   url,
		dao.TableQrcodes.Columns().CreatedAt:   gtime.Now(),
		dao.TableQrcodes.Columns().UpdatedAt:   gtime.Now(),
	}

	// 插入数据库
	result, err := dao.TableQrcodes.Ctx(ctx).Data(data).Insert()
	if err != nil {
		return nil, gerror.Wrap(err, "保存二维码信息失败")
	}

	// 获取插入的ID
	id, err := result.LastInsertId()
	if err != nil {
		return nil, gerror.Wrap(err, "获取二维码ID失败")
	}

	// 6. 返回结果
	return &qrcodev1.TableQrcode{
		ID:          id,
		TableNumber: req.TableNumber,
		QrcodeURL:   url,
		CreatedAt:   gtime.Now().String(),
	}, nil
}

// GetTableQrcodeList 获取桌号二维码列表
func (s *qrcodeService) GetTableQrcodeList(ctx context.Context, req *qrcodev1.TableQrcodeListReq) ([]qrcodev1.TableQrcode, int, error) {
	// 1. 构建查询条件
	m := dao.TableQrcodes.Ctx(ctx)

	// 1.1 桌号模糊搜索
	if req.TableNumber != "" {
		m = m.WhereLike(dao.TableQrcodes.Columns().TableNumber, "%"+req.TableNumber+"%")
	}

	// 2. 查询总数
	total, err := m.Count()
	if err != nil {
		return nil, 0, err
	}

	// 3. 分页参数处理
	page := req.Page
	if page < 1 {
		page = 1
	}
	limit := req.Limit
	if limit < 1 {
		limit = 10
	}

	// 4. 查询数据
	var tableQrcodes []*entity.TableQrcodes
	err = m.Page(page, limit).
		Order(dao.TableQrcodes.Columns().Id + " DESC").
		Scan(&tableQrcodes)
	if err != nil {
		return nil, 0, err
	}

	// 5. 转换为API响应格式
	result := make([]qrcodev1.TableQrcode, len(tableQrcodes))
	for i, qrcode := range tableQrcodes {
		result[i] = qrcodev1.TableQrcode{
			ID:          int64(qrcode.Id),
			TableNumber: qrcode.TableNumber,
			QrcodeURL:   qrcode.QrcodeUrl,
			CreatedAt:   qrcode.CreatedAt.String(),
		}
	}

	return result, total, nil
}

// DeleteTableQrcode 删除桌号二维码
func (s *qrcodeService) DeleteTableQrcode(ctx context.Context, id int64) error {
	// 1. 检查桌号是否存在
	var tableQrcode *entity.TableQrcodes
	err := dao.TableQrcodes.Ctx(ctx).
		Where(dao.TableQrcodes.Columns().Id, id).
		Scan(&tableQrcode)
	if err != nil {
		return err
	}
	if tableQrcode == nil {
		return gerror.New("桌号不存在")
	}

	// 2. 检查桌号是否有关联的订单
	count, err := dao.Orders.Ctx(ctx).
		Where(dao.Orders.Columns().TableQrcodeId, id).
		Count()
	if err != nil {
		return err
	}
	if count > 0 {
		return gerror.New("该桌号已有关联订单，无法删除")
	}

	// 3. 删除桌号二维码记录
	_, err = dao.TableQrcodes.Ctx(ctx).
		Where(dao.TableQrcodes.Columns().Id, id).
		Delete()
	if err != nil {
		return gerror.Wrap(err, "删除桌号失败")
	}

	// 4. 可以选择删除MinIO中的二维码图片，但由于图片URL可能被缓存，这里暂不删除
	// 如果需要删除，可以通过tableQrcode.QrcodeUrl解析出文件名，然后调用MinIO的删除方法

	return nil
}
