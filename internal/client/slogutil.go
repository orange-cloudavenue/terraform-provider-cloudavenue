/*
 * SPDX-FileCopyrightText: Copyright (c) 2025 Orange
 * SPDX-License-Identifier: Mozilla Public License 2.0
 *
 * This software is distributed under the MPL-2.0 license.
 * the text of which is available at https://www.mozilla.org/en-US/MPL/2.0/
 * or see the "LICENSE" file for more details.
 */

package client

import (
	"context"
	"log/slog"

	"github.com/hashicorp/terraform-plugin-log/tflog"
)

func newTflogHandler() slog.Handler {
	return &slogHandler{}
}

type slogHandler struct {
	attrs  []slog.Attr
	groups []string
}

var _ slog.Handler = (*slogHandler)(nil)

func (*slogHandler) Enabled(context.Context, slog.Level) bool {
	return true
}

func (h *slogHandler) Handle(ctx context.Context, record slog.Record) error {
	switch record.Level {
	case slog.LevelDebug:
		tflog.Debug(ctx, record.Message, h.fields(record))
	case slog.LevelInfo:
		tflog.Info(ctx, record.Message, h.fields(record))
	case slog.LevelWarn:
		tflog.Warn(ctx, record.Message, h.fields(record))
	case slog.LevelError:
		tflog.Error(ctx, record.Message, h.fields(record))
	default:
		tflog.Info(ctx, record.Message, h.fields(record))
	}
	return nil
}

func (h *slogHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &slogHandler{
		attrs:  append(h.attrs, attrs...),
		groups: h.groups,
	}
}

func (h *slogHandler) WithGroup(name string) slog.Handler {
	if name == "" {
		return h
	}
	return &slogHandler{
		attrs:  h.attrs,
		groups: append(h.groups, name),
	}
}

func (h *slogHandler) fields(record slog.Record) map[string]any {
	root := make(map[string]any, len(h.attrs)+record.NumAttrs())

	fields := root
	for _, name := range h.groups {
		nested := make(map[string]any, len(h.attrs)+record.NumAttrs())
		fields[name] = nested
		fields = nested
	}

	addAttrsToMap(h.attrs, fields)

	record.Attrs(func(attr slog.Attr) bool {
		addAttrToMap(attr, fields)
		return true
	})

	return root
}

func addAttrsToMap(attrs []slog.Attr, fields map[string]any) {
	for _, a := range attrs {
		addAttrToMap(a, fields)
	}
}

func addAttrToMap(attr slog.Attr, fields map[string]any) {
	if attr.Equal(slog.Attr{}) {
		return
	}

	val := attr.Value.Resolve()

	if val.Kind() == slog.KindGroup {
		attrs := val.Group()
		if len(attrs) == 0 {
			return
		}

		if attr.Key == "" {
			addAttrsToMap(attrs, fields)
			return
		}

		group := make(map[string]any, len(attrs))
		addAttrsToMap(attrs, group)
		fields[attr.Key] = group

		return
	}

	fields[attr.Key] = val.Any()
}
