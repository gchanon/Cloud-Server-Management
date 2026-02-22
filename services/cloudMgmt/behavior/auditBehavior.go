package behavior

import (
	"os"
	"time"

	"github.com/golf/cloudmgmt/services/cloudMgmt/model"
	"github.com/golf/cloudmgmt/services/cloudMgmt/utility"
)

type AuditTrail struct {
	auditData map[int64]map[string]*model.AuditTrailModel
}

func NewAuditTrailBehavior() *AuditTrail {
	return &AuditTrail{
		auditData: make(map[int64]map[string]*model.AuditTrailModel),
	}
}

func (audit *AuditTrail) Create(auditTrail *model.AuditTrailModel) {
	if _, exists := audit.auditData[auditTrail.UserID]; !exists {
		audit.auditData[auditTrail.UserID] = make(map[string]*model.AuditTrailModel)
	}

	audit.auditData[auditTrail.UserID][utility.NewChronoSequence(time.Now(), os.Getpid())] = auditTrail

}
