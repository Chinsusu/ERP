package entity_test

import (
	"testing"

	"github.com/erp-cosmetics/manufacturing-service/internal/domain/entity"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestBOM_Encryption(t *testing.T) {
	key := []byte("thisis32bytekeyforaesgcmtesting!") // 32 bytes for AES-256
	formula := &entity.FormulaDetails{
		ProcessingSteps: []string{"Step 1: Mix", "Step 2: Heat"},
		CriticalParameters: map[string]string{
			"Temperature": "80C",
			"MixSpeed":    "100rpm",
		},
		Notes: "Confidential formula",
	}

	t.Run("Encrypt and Decrypt Success", func(t *testing.T) {
		encrypted, err := entity.EncryptFormula(formula, key)
		assert.NoError(t, err)
		assert.NotEmpty(t, encrypted)
		assert.NotEqual(t, formula.Notes, string(encrypted)) // Should not be readable

		decrypted, err := entity.DecryptFormula(encrypted, key)
		assert.NoError(t, err)
		assert.Equal(t, formula.Notes, decrypted.Notes)
		assert.Equal(t, formula.ProcessingSteps, decrypted.ProcessingSteps)
	})

	t.Run("Decrypt with Wrong Key Fails", func(t *testing.T) {
		encrypted, _ := entity.EncryptFormula(formula, key)
		wrongKey := []byte("wrongkeywrongkeywrongkeywrongkey")
		
		decrypted, err := entity.DecryptFormula(encrypted, wrongKey)
		assert.Error(t, err)
		assert.Nil(t, decrypted)
	})

	t.Run("Decrypt Invalid Ciphertext Fails", func(t *testing.T) {
		invalidData := []byte("invalid-encrypted-data")
		decrypted, err := entity.DecryptFormula(invalidData, key)
		assert.Error(t, err)
		assert.Nil(t, decrypted)
	})
}

func TestBOM_StatusTransitions(t *testing.T) {
	bom := &entity.BOM{
		Status: entity.BOMStatusDraft,
	}

	t.Run("Draft to Pending Approval", func(t *testing.T) {
		err := bom.Submit()
		assert.NoError(t, err)
		assert.Equal(t, entity.BOMStatusPendingApproval, bom.Status)
	})

	t.Run("Pending Approval to Approved", func(t *testing.T) {
		approverID := uuid.New()
		err := bom.Approve(approverID)
		assert.NoError(t, err)
		assert.Equal(t, entity.BOMStatusApproved, bom.Status)
		assert.Equal(t, &approverID, bom.ApprovedBy)
		assert.NotNil(t, bom.ApprovedAt)
	})

	t.Run("Approve Draft Fails", func(t *testing.T) {
		bom.Status = entity.BOMStatusDraft
		err := bom.Approve(uuid.New())
		assert.Error(t, err)
	})
}
