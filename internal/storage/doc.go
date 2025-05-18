// Package storage provides implementations for URL shortener storage backends.
//
// It includes three ready-to-use implementations:
//   - In-memory storage (AppMemStorage) - for testing and ephemeral data
//   - File-based storage (FileStorage) - for single-instance persistence
//   - Database storage (DatabaseStorage) - for production PostgreSQL deployments
package storage
