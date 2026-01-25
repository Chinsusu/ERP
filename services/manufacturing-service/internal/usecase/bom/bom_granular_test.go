package bom_test

import (
	"context"
	"testing"

	"github.com/erp-cosmetics/manufacturing-service/internal/domain/entity"
	"github.com/erp-cosmetics/manufacturing-service/internal/testmocks"
	"github.com/erp-cosmetics/manufacturing-service/internal/usecase/bom"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestGetBOMUseCase_Execute_GranularPermissions(t *testing.T) {
	// Arrange
	ctx := context.Background()
	repo := new(testmocks.MockBOMRepository)
	key := []byte("thisis32bytekeyforaesgcmtesting!")
	
	uc := bom.NewGetBOMUseCase(repo, key)
	bomID := uuid.New()
	
	formula := &entity.FormulaDetails{Notes: "Confidential Process"}
	encrypted, _ := entity.EncryptFormula(formula, key)
	
	targetBOM := &entity.BOM{
		ID:             bomID,
		FormulaDetails: encrypted,
		Items: []entity.BOMLineItem{
			{Quantity: 100},
		},
	}

	repo.On("GetByID", ctx, bomID).Return(targetBOM, nil)

	tests := []struct {
		name           string
		canViewFormula bool
		expectedNotes  string
		expectNil      bool
	}{
		{
			name:           "User with formula access",
			canViewFormula: true,
			expectedNotes:  "Confidential Process",
			expectNil:      false,
		},
		{
			name:           "User without formula access",
			canViewFormula: false,
			expectedNotes:  "",
			expectNil:      true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res, err := uc.Execute(ctx, bomID, tt.canViewFormula)
			assert.NoError(t, err)
			if tt.expectNil {
				assert.Nil(t, res.FormulaDetails)
			} else {
				assert.NotNil(t, res.FormulaDetails)
				assert.Equal(t, tt.expectedNotes, res.FormulaDetails.Notes)
			}
		})
	}
}

func TestBOM_ReEncryption(t *testing.T) {
	oldKey := []byte("oldkeyoldkeyoldkeyoldkeyoldkey32")
	newKey := []byte("newkeynewkeynewkeynewkeynewkey32")
	formula := &entity.FormulaDetails{Notes: "Top Secret"}

	// 1. Encrypt with old key
	encrypted, _ := entity.EncryptFormula(formula, oldKey)
	
	// 2. Decrypt with old key
	decrypted, err := entity.DecryptFormula(encrypted, oldKey)
	assert.NoError(t, err)
	assert.Equal(t, "Top Secret", decrypted.Notes)

	// 3. Re-encrypt with new key
	reEncrypted, err := entity.EncryptFormula(decrypted, newKey)
	assert.NoError(t, err)
	assert.NotEmpty(t, reEncrypted)
	assert.NotEqual(t, encrypted, reEncrypted)

	// 4. Verify with new key
	finalDecrypted, err := entity.DecryptFormula(reEncrypted, newKey)
	assert.NoError(t, err)
	assert.Equal(t, "Top Secret", finalDecrypted.Notes)
}
