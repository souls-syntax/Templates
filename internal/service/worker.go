package service

import (
	"context"
	"fmt"
	"time"

	"github.com/souls-syntax/Templates/internal/cache"
	"github.com/souls-syntax/Templates/internal/database"
	"github.com/souls-syntax/Templates/internal/models"
)

type AnalysisJob struct {
	Hash					 string
	QueryText			 string
}
