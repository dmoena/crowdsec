package controllers

import (
	"context"
	"fmt"
	"net/http"

	"github.com/crowdsecurity/crowdsec/pkg/database/ent"
	"github.com/crowdsecurity/crowdsec/pkg/database/ent/machine"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

type CreateMachineInput struct {
	MachineId string `json:"machine_id" binding:"required"`
	Password  string `json:"password" binding:"required"`
	IpAddress string `json:"ip_address" binding:"required"`
}

func QueryMachine(ctx context.Context, client *ent.Client, machineId string) (*ent.Machine, error) {
	machine, err := client.Machine.
		Query().
		Where(machine.MachineIdEQ(machineId)).
		Only(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed querying user: %v", err)
	}
	return machine, nil
}

func (c *Controller) CreateMachine(gctx *gin.Context) {
	var input CreateMachineInput
	if err := gctx.ShouldBindJSON(&input); err != nil {
		gctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	machine, err := c.Client.Machine.
		Create().
		SetMachineId(input.MachineId).
		SetPassword(input.Password).
		SetIpAddress(input.IpAddress).
		Save(c.Ectx)

	if err != nil {
		log.Errorf("failed creating machine: %v", err)
		gctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed creating machine"})
		return
	}
	gctx.JSON(http.StatusOK, gin.H{"data": machine})
	return
}