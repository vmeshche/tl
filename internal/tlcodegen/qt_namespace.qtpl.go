// Copyright 2022 V Kontakte LLC
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

// Code generated by qtc from "qt_namespace.qtpl". DO NOT EDIT.
// See https://github.com/valyala/quicktemplate for details.

package tlcodegen

import (
	"fmt"

	// will generate type aliases anyway and RPC code if appropriate flag is set

	qtio422016 "io"

	qt422016 "github.com/valyala/quicktemplate"
)

var (
	_ = qtio422016.Copy
	_ = qt422016.AcquireByteBuffer
)

func (gen *Gen2) streamgenerateNamespacesCode(qw422016 *qt422016.Writer, anyTypeAlias bool, anyFunctions bool, name string, namespace *Namespace, sortedImports []string, directImports *DirectImports) {
	qw422016.N().S(HeaderComment)
	qw422016.N().S(`
package `)
	qw422016.N().S(gen.GlobalPackageName + name)
	qw422016.N().S(`

import (
`)
	if gen.options.GenerateRPCCode && anyFunctions {
		qw422016.N().S(`    "context"

    `)
		qw422016.N().Q(gen.options.BasicRPCPath)
		qw422016.N().S(`
    `)
		qw422016.N().Q(gen.BasicPackageNameFull)
		qw422016.N().S(`
    "`)
		qw422016.N().S(gen.options.TLPackageNameFull)
		qw422016.N().S(`/internal"
`)
	}
	for _, wr := range sortedImports {
		qw422016.N().S(`    "`)
		qw422016.N().S(gen.options.TLPackageNameFull)
		qw422016.N().S(`/`)
		qw422016.N().S(wr)
		qw422016.N().S(`"
`)
	}
	qw422016.N().S(`)

`)
	ourTypes := map[*TypeRWWrapper]struct{}{}

	qw422016.N().S(`
`)
	streamtypesAlias(qw422016, anyTypeAlias, name, namespace.types, directImports, ourTypes)
	qw422016.N().S(`

`)
	if !gen.options.GenerateRPCCode || !anyFunctions {
		return
	}
	qw422016.N().S(`
type Client struct {
    Client  *rpc.Client
    Network string // should be either "tcp4" or "unix"
    Address string
    ActorID uint64 // should be non-zero when using rpc-proxy
}

`)
	streamwriteClientsCode(qw422016, gen.GlobalPackageName, namespace.types, directImports, ourTypes)
	qw422016.N().S(`

type Handler struct {
    `)
	streamhandlerStructs(qw422016, gen.GlobalPackageName, name, namespace.types, directImports, ourTypes)
	qw422016.N().S(` }

func (h *Handler) Handle(ctx context.Context, hctx *rpc.HandlerContext) (err error) {
    `)
	streamhandleRequest(qw422016, namespace.types, directImports)
	qw422016.N().S(` return rpc.ErrNoHandler
}
`)
}

func (gen *Gen2) writegenerateNamespacesCode(qq422016 qtio422016.Writer, anyTypeAlias bool, anyFunctions bool, name string, namespace *Namespace, sortedImports []string, directImports *DirectImports) {
	qw422016 := qt422016.AcquireWriter(qq422016)
	gen.streamgenerateNamespacesCode(qw422016, anyTypeAlias, anyFunctions, name, namespace, sortedImports, directImports)
	qt422016.ReleaseWriter(qw422016)
}

func (gen *Gen2) generateNamespacesCode(anyTypeAlias bool, anyFunctions bool, name string, namespace *Namespace, sortedImports []string, directImports *DirectImports) string {
	qb422016 := qt422016.AcquireByteBuffer()
	gen.writegenerateNamespacesCode(qb422016, anyTypeAlias, anyFunctions, name, namespace, sortedImports, directImports)
	qs422016 := string(qb422016.B)
	qt422016.ReleaseByteBuffer(qb422016)
	return qs422016
}

func streamtypesAlias(qw422016 *qt422016.Writer, anyTypeAlias bool, namespace string, types []*TypeRWWrapper, directImports *DirectImports, ourTypes map[*TypeRWWrapper]struct{}) {
	if anyTypeAlias {
		qw422016.N().S(`type(
`)
		for _, wr := range types {
			if wr.ShouldWriteTypeAlias() {
				ourTypes[wr] = struct{}{}

				qw422016.N().S(`    `)
				qw422016.N().S(wr.TypeString2(false, directImports, nil, true, true))
				qw422016.N().S(` = `)
				qw422016.N().S(wr.TypeString2(false, directImports, nil, false, true))
				qw422016.N().S(`
`)
				if wr.wantsBytesVersion && wr.hasBytesVersion {
					qw422016.N().S(`    `)
					qw422016.N().S(wr.TypeString2(true, directImports, nil, true, true))
					qw422016.N().S(` = `)
					qw422016.N().S(wr.TypeString2(true, directImports, nil, false, true))
					qw422016.N().S(`
`)
				}
			}
		}
		qw422016.N().S(`)
`)
	}
	for _, wr := range types {
		if wr.ShouldWriteEnumElementAlias() {
			_, ourUnionParentType := ourTypes[wr.unionParent]
			typeString := wr.TypeString2(false, directImports, nil, true, true)

			qw422016.N().S(`func `)
			qw422016.N().S(typeString)
			qw422016.N().S(`() `)
			qw422016.N().S(wr.unionParent.TypeString2(false, directImports, nil, ourUnionParentType, false))
			qw422016.N().S(` { return `)
			qw422016.N().S(wr.TypeString2(false, directImports, nil, false, true))
			qw422016.N().S(`() }
`)
		}
	}
}

func writetypesAlias(qq422016 qtio422016.Writer, anyTypeAlias bool, namespace string, types []*TypeRWWrapper, directImports *DirectImports, ourTypes map[*TypeRWWrapper]struct{}) {
	qw422016 := qt422016.AcquireWriter(qq422016)
	streamtypesAlias(qw422016, anyTypeAlias, namespace, types, directImports, ourTypes)
	qt422016.ReleaseWriter(qw422016)
}

func typesAlias(anyTypeAlias bool, namespace string, types []*TypeRWWrapper, directImports *DirectImports, ourTypes map[*TypeRWWrapper]struct{}) string {
	qb422016 := qt422016.AcquireByteBuffer()
	writetypesAlias(qb422016, anyTypeAlias, namespace, types, directImports, ourTypes)
	qs422016 := string(qb422016.B)
	qt422016.ReleaseByteBuffer(qb422016)
	return qs422016
}

func streamwriteClientsCode(qw422016 *qt422016.Writer, shortPackageName string, types []*TypeRWWrapper, directImports *DirectImports, ourTypes map[*TypeRWWrapper]struct{}) {
	for _, wr := range types {
		if wr.wantsBytesVersion && wr.hasBytesVersion {
			streamwriteClientCode(qw422016, true, shortPackageName, wr, directImports, ourTypes)
			qw422016.N().S(`
`)
		}
		streamwriteClientCode(qw422016, false, shortPackageName, wr, directImports, ourTypes)
		qw422016.N().S(`
`)
	}
}

func writewriteClientsCode(qq422016 qtio422016.Writer, shortPackageName string, types []*TypeRWWrapper, directImports *DirectImports, ourTypes map[*TypeRWWrapper]struct{}) {
	qw422016 := qt422016.AcquireWriter(qq422016)
	streamwriteClientsCode(qw422016, shortPackageName, types, directImports, ourTypes)
	qt422016.ReleaseWriter(qw422016)
}

func writeClientsCode(shortPackageName string, types []*TypeRWWrapper, directImports *DirectImports, ourTypes map[*TypeRWWrapper]struct{}) string {
	qb422016 := qt422016.AcquireByteBuffer()
	writewriteClientsCode(qb422016, shortPackageName, types, directImports, ourTypes)
	qs422016 := string(qb422016.B)
	qt422016.ReleaseByteBuffer(qb422016)
	return qs422016
}

func streamwriteClientCode(qw422016 *qt422016.Writer, bytesVersion bool, shortPackageName string, wr *TypeRWWrapper, directImports *DirectImports, ourTypes map[*TypeRWWrapper]struct{}) {
	fun, ok := wr.trw.(*TypeRWStruct)

	if !ok || fun.ResultType == nil {
		return
	}
	_, ourResultType := ourTypes[fun.ResultType]
	ret := fun.ResultType.TypeString2(bytesVersion, directImports, nil, ourResultType, false)
	typeString := wr.TypeString2(bytesVersion, directImports, nil, true, true)
	tlName := wr.tlName.String()

	qw422016.N().S(`func (c *Client) `)
	qw422016.N().S(typeString)
	qw422016.N().S(`(ctx context.Context, args `)
	qw422016.N().S(typeString)
	qw422016.N().S(`, extra *rpc.InvokeReqExtra, ret *`)
	qw422016.N().S(ret)
	qw422016.N().S(`) (err error) {
    req := c.Client.GetRequest()
    req.ActorID = c.ActorID
    if extra != nil {
        req.Extra = *extra
    }
    req.Body, err = args.WriteBoxed(req.Body)
    if err != nil {
        return internal.ErrorClientWrite("`)
	qw422016.N().S(tlName)
	qw422016.N().S(`", err)
    }
    resp, err := c.Client.Do(ctx, c.Network, c.Address, req)
    if err != nil {
        return internal.ErrorClientDo("`)
	qw422016.N().S(tlName)
	qw422016.N().S(`", c.Network, c.ActorID, c.Address, err)
    }
    defer c.Client.PutResponse(resp)
    if ret != nil {
        if _, err = args.ReadResult(resp.Body, ret); err != nil {
            return internal.ErrorClientReadResult("`)
	qw422016.N().S(tlName)
	qw422016.N().S(`", c.Network, c.ActorID, c.Address, err)
        }
    }
    return nil
}
`)
}

func writewriteClientCode(qq422016 qtio422016.Writer, bytesVersion bool, shortPackageName string, wr *TypeRWWrapper, directImports *DirectImports, ourTypes map[*TypeRWWrapper]struct{}) {
	qw422016 := qt422016.AcquireWriter(qq422016)
	streamwriteClientCode(qw422016, bytesVersion, shortPackageName, wr, directImports, ourTypes)
	qt422016.ReleaseWriter(qw422016)
}

func writeClientCode(bytesVersion bool, shortPackageName string, wr *TypeRWWrapper, directImports *DirectImports, ourTypes map[*TypeRWWrapper]struct{}) string {
	qb422016 := qt422016.AcquireByteBuffer()
	writewriteClientCode(qb422016, bytesVersion, shortPackageName, wr, directImports, ourTypes)
	qs422016 := string(qb422016.B)
	qt422016.ReleaseByteBuffer(qb422016)
	return qs422016
}

func streamhandlerStructs(qw422016 *qt422016.Writer, shortPackageName string, name string, types []*TypeRWWrapper, directImports *DirectImports, ourTypes map[*TypeRWWrapper]struct{}) {
	for _, wr := range types {
		if fun, ok := wr.trw.(*TypeRWStruct); ok && fun.ResultType != nil {
			tlName := wr.tlName.String()
			_, ourResultType := ourTypes[fun.ResultType]
			ret := fun.ResultType.TypeString2(false, directImports, nil, ourResultType, false)
			funcTypeString := wr.TypeString2(false, directImports, nil, true, true)

			qw422016.N().S(funcTypeString)
			qw422016.N().S(` func(ctx context.Context, args `)
			qw422016.N().S(funcTypeString)
			qw422016.N().S(`) (`)
			qw422016.N().S(ret)
			qw422016.N().S(`, error) // `)
			qw422016.N().S(tlName)
			qw422016.N().S(`
`)
		}
	}
	qw422016.N().S(`
`)
	for _, wr := range types {
		if fun, ok := wr.trw.(*TypeRWStruct); ok && fun.ResultType != nil {
			tlName := wr.tlName.String()
			funcTypeString := wr.TypeString2(false, directImports, nil, true, true)

			qw422016.N().S(`Raw`)
			qw422016.N().S(funcTypeString)
			qw422016.N().S(` func(ctx context.Context, hctx *rpc.HandlerContext) error // `)
			qw422016.N().S(tlName)
			qw422016.N().S(`
`)
		}
	}
}

func writehandlerStructs(qq422016 qtio422016.Writer, shortPackageName string, name string, types []*TypeRWWrapper, directImports *DirectImports, ourTypes map[*TypeRWWrapper]struct{}) {
	qw422016 := qt422016.AcquireWriter(qq422016)
	streamhandlerStructs(qw422016, shortPackageName, name, types, directImports, ourTypes)
	qt422016.ReleaseWriter(qw422016)
}

func handlerStructs(shortPackageName string, name string, types []*TypeRWWrapper, directImports *DirectImports, ourTypes map[*TypeRWWrapper]struct{}) string {
	qb422016 := qt422016.AcquireByteBuffer()
	writehandlerStructs(qb422016, shortPackageName, name, types, directImports, ourTypes)
	qs422016 := string(qb422016.B)
	qt422016.ReleaseByteBuffer(qb422016)
	return qs422016
}

func streamhandleRequest(qw422016 *qt422016.Writer, types []*TypeRWWrapper, directImports *DirectImports) {
	totalFuns := 0
	for _, wr := range types {
		if fun, ok := wr.trw.(*TypeRWStruct); ok && fun.ResultType != nil {
			totalFuns++
		}
	}

	if totalFuns == 0 {
		return
	}
	qw422016.N().S(`tag, r, _ := basictl.NatReadTag(hctx.Request) // keep hctx.Request intact for handler chaining
switch tag {
`)
	// TODO - we have to skip bytes in bytes.Buffer for Raw call
	// TODO - check that no bytes remains after reading

	for _, wr := range types {
		if fun, ok := wr.trw.(*TypeRWStruct); ok && fun.ResultType != nil {
			tlTag := fmt.Sprintf("%#08x", wr.tlTag)
			funcTypeString := wr.TypeString2(false, directImports, nil, true, true)
			tlName := wr.tlName.String()

			qw422016.N().S(`case `)
			qw422016.N().S(tlTag)
			qw422016.N().S(`: // `)
			qw422016.N().S(tlName)
			qw422016.N().S(`
    if h.Raw`)
			qw422016.N().S(funcTypeString)
			qw422016.N().S(` != nil {
        hctx.Request = r
        err = h.Raw`)
			qw422016.N().S(funcTypeString)
			qw422016.N().S(`(ctx, hctx)
        if rpc.IsHijackedResponse(err) {
            return err
        }
        if err != nil {
            return internal.ErrorServerHandle("`)
			qw422016.N().S(tlName)
			qw422016.N().S(`", err)
        }
        return nil
    }
    if h.`)
			qw422016.N().S(funcTypeString)
			qw422016.N().S(` != nil {
        var args `)
			qw422016.N().S(funcTypeString)
			qw422016.N().S(`
        if _, err = args.Read(r); err != nil {
            return internal.ErrorServerRead("`)
			qw422016.N().S(tlName)
			qw422016.N().S(`", err)
        }
        ctx = hctx.WithContext(ctx)
        ret, err := h.`)
			qw422016.N().S(funcTypeString)
			qw422016.N().S(`(ctx, args)
        if rpc.IsHijackedResponse(err)  {
            return err
        }
        if err != nil {
            return internal.ErrorServerHandle("`)
			qw422016.N().S(tlName)
			qw422016.N().S(`", err)
        }
        if hctx.Response, err = args.WriteResult(hctx.Response, ret); err != nil {
            return internal.ErrorServerWriteResult("`)
			qw422016.N().S(tlName)
			qw422016.N().S(`", err)
        }
        return nil
    }
`)
		}
	}
	qw422016.N().S(`}
`)
}

func writehandleRequest(qq422016 qtio422016.Writer, types []*TypeRWWrapper, directImports *DirectImports) {
	qw422016 := qt422016.AcquireWriter(qq422016)
	streamhandleRequest(qw422016, types, directImports)
	qt422016.ReleaseWriter(qw422016)
}

func handleRequest(types []*TypeRWWrapper, directImports *DirectImports) string {
	qb422016 := qt422016.AcquireByteBuffer()
	writehandleRequest(qb422016, types, directImports)
	qs422016 := string(qb422016.B)
	qt422016.ReleaseByteBuffer(qb422016)
	return qs422016
}
