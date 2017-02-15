package dao

type Dao struct {
	Redis *Redis
	HBase *HBase
	Mysql *Mysql
}

func NewDao() (dao *Dao) {
	redis := NewRedis()
	h := NewHBase()
	mysql := NewMysql()
	dao = &Dao{
		Redis: redis,
		HBase: h,
		Mysql: mysql,
	}
	return
}
