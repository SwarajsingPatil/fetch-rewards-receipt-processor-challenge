package store

import (
    "receipt-processor/internal/models"
    "sync"
)

type MemoryStore struct {
    receipts map[string]*models.Receipt
    points   map[string]int64
    mu       sync.RWMutex
}

func NewMemoryStore() *MemoryStore {
    return &MemoryStore{
        receipts: make(map[string]*models.Receipt),
        points:   make(map[string]int64),
    }
}

func (s *MemoryStore) SaveReceipt(id string, receipt *models.Receipt, points int64) {
    s.mu.Lock()
    defer s.mu.Unlock()
    s.receipts[id] = receipt
    s.points[id] = points
}

func (s *MemoryStore) GetPoints(id string) (int64, bool) {
    s.mu.RLock()
    defer s.mu.RUnlock()
    points, exists := s.points[id]
    return points, exists
}