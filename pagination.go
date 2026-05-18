package hooksniff

import (
	"context"
)

// ListResponse is a generic interface for paginated list responses.
type ListResponse[T any] interface {
	GetData() []T
	GetDone() bool
	GetIterator() *string
}

// Paginator provides automatic cursor-based pagination.
type Paginator[T any] struct {
	fetchPage func(ctx context.Context, iterator *string) (ListResponse[T], error)
	ctx       context.Context
	current   ListResponse[T]
	index     int
	done      bool
	iterator  *string
	err       error
}

// NewPaginator creates a new paginator that auto-fetches pages.
//
// Usage:
//
//	paginator := hooksniff.NewPaginator(ctx, func(ctx context.Context, iter *string) (ListResponse[T], error) {
//	    return client.Message.List(ctx, &MessageListOptions{Limit: lo(100), Iterator: iter})
//	})
//	for paginator.Next() {
//	    item := paginator.Value()
//	    // process item
//	}
//	if err := paginator.Err(); err != nil {
//	    // handle error
//	}
func NewPaginator[T any](ctx context.Context, fetchPage func(ctx context.Context, iterator *string) (ListResponse[T], error)) *Paginator[T] {
	return &Paginator[T]{
		fetchPage: fetchPage,
		ctx:       ctx,
	}
}

// Next advances to the next item. Returns false when done or on error.
func (p *Paginator[T]) Next() bool {
	if p.done {
		return false
	}

	// Fetch first page if needed
	if p.current == nil {
		page, err := p.fetchPage(p.ctx, p.iterator)
		if err != nil {
			p.err = err
			return false
		}
		p.current = page
		p.index = 0
	}

	// Return items from current page
	if p.index < len(p.current.GetData()) {
		return true
	}

	// Move to next page
	if !p.current.GetDone() && p.current.GetIterator() != nil {
		p.iterator = p.current.GetIterator()
		page, err := p.fetchPage(p.ctx, p.iterator)
		if err != nil {
			p.err = err
			return false
		}
		p.current = page
		p.index = 0
		return len(p.current.GetData()) > 0
	}

	p.done = true
	return false
}

// Value returns the current item. Call after Next() returns true.
func (p *Paginator[T]) Value() T {
	item := p.current.GetData()[p.index]
	p.index++
	return item
}

// Err returns any error that occurred during pagination.
func (p *Paginator[T]) Err() error {
	return p.err
}

// All collects all remaining items into a slice.
func (p *Paginator[T]) All() ([]T, error) {
	var items []T
	for p.Next() {
		items = append(items, p.Value())
	}
	return items, p.Err()
}

// PagePaginator provides page-level pagination (returns full pages).
type PagePaginator[T any] struct {
	fetchPage func(ctx context.Context, iterator *string) (ListResponse[T], error)
	ctx       context.Context
	iterator  *string
	done      bool
	err       error
}

// NewPagePaginator creates a paginator that yields full pages.
func NewPagePaginator[T any](ctx context.Context, fetchPage func(ctx context.Context, iterator *string) (ListResponse[T], error)) *PagePaginator[T] {
	return &PagePaginator[T]{
		fetchPage: fetchPage,
		ctx:       ctx,
	}
}

// NextPage fetches the next page. Returns nil when done.
func (p *PagePaginator[T]) NextPage() (ListResponse[T], error) {
	if p.done {
		return nil, nil
	}

	page, err := p.fetchPage(p.ctx, p.iterator)
	if err != nil {
		p.err = err
		return nil, err
	}

	if page.GetDone() || page.GetIterator() == nil {
		p.done = true
	} else {
		p.iterator = page.GetIterator()
	}

	return page, nil
}

// Helper to create *string from string
func strPtr(s string) *string {
	return &s
}
