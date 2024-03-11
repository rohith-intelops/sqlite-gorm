package daos

import (
	"errors"
	"github.com/rohith-intelops/socialmedia/user-service/pkg/rest/server/daos/clients/sqls"
	"github.com/rohith-intelops/socialmedia/user-service/pkg/rest/server/models"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type UserDao struct {
	db *gorm.DB
}

func NewUserDao() (*UserDao, error) {
	sqlClient, err := sqls.InitGORMSQLiteDB()
	if err != nil {
		return nil, err
	}
	err = sqlClient.DB.AutoMigrate(models.User{})
	if err != nil {
		return nil, err
	}
	return &UserDao{
		db: sqlClient.DB,
	}, nil
}

func (userDao *UserDao) CreateUser(m *models.User) (*models.User, error) {
	if err := userDao.db.Create(&m).Error; err != nil {
		log.Debugf("failed to create user: %v", err)
		return nil, err
	}

	log.Debugf("user created")
	return m, nil
}

func (userDao *UserDao) GetUser(id int64) (*models.User, error) {
	var m *models.User
	if err := userDao.db.Where("id = ?", id).First(&m).Error; err != nil {
		log.Debugf("failed to get user: %v", err)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, sqls.ErrNotExists
		}
		return nil, err
	}
	log.Debugf("user retrieved")
	return m, nil
	
}

func (userDao *UserDao) UpdateUser(id int64, m *models.User) (*models.User, error) {
	if id == 0 {
		return nil, errors.New("invalid user ID")
	}
	if id != m.Id {
		return nil, errors.New("id and payload don't match")
	}

	var user *models.User
	if err := userDao.db.Where("id = ?", id).First(&user).Error; err != nil {
		log.Debugf("failed to find user for update: %v", err)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, sqls.ErrNotExists
		}
		return nil, err
	}

	if err := userDao.db.Save(&m).Error; err != nil {
		log.Debugf("failed to update user: %v", err)
		return nil, err
	}
	log.Debugf("user updated")
	return m, nil
}

func (userDao *UserDao) DeleteUser(id int64) error {
	var m *models.User
	if err := userDao.db.Where("id = ?", id).Delete(&m).Error; err != nil {
		log.Debugf("failed to delete user: %v", err)
		return err
	}

	log.Debugf("user deleted")
	return nil
}
