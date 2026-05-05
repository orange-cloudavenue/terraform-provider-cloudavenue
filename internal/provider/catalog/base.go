/*
 * SPDX-FileCopyrightText: Copyright (c) 2026 Orange
 * SPDX-License-Identifier: Mozilla Public License 2.0
 *
 * This software is distributed under the MPL-2.0 license.
 * the text of which is available at https://www.mozilla.org/en-US/MPL/2.0/
 * or see the "LICENSE" file for more details.
 */

/*
 * SPDX-FileCopyrightText: Copyright (c) 2025 Orange
 * SPDX-License-Identifier: Mozilla Public License 2.0
 *
 * This software is distributed under the MPL-2.0 license.
 * the text of which is available at https://www.mozilla.org/en-US/MPL/2.0/
 * or see the "LICENSE" file for more details.
 */

package catalog

import (
	"github.com/vmware/go-vcloud-director/v2/govcd"
)

const (
	categoryName = "catalog"

	catalogID   = "catalog_id"
	catalogName = "catalog_name"

	// Attribute keys.
	createdAt      = "created_at"
	description    = "description"
	isPublished    = "is_published"
	name           = "name"
	ownerName      = "owner_name"
	storageProfile = "storage_profile"

	// Attribute keys (additional).
	catalogsAttr                = "catalogs"
	isCached                    = "is_cached"
	isISO                       = "is_iso"
	isLocal                     = "is_local"
	isShared                    = "is_shared"
	mediaItemList               = "media_item_list"
	numberOfMedia               = "number_of_media"
	preserveIdentityInformation = "preserve_identity_information"
	sharedWithEveryone          = "shared_with_everyone"
	size                        = "size"
	status                      = "status"
	templateID                  = "template_id"
	templateName                = "template_name"

	// Attribute descriptions.
	catalogIDDescription   = "The ID of the catalog."
	catalogNameDescription = "The name of the catalog."

	// Shared MarkdownDescriptions.
	descCatalogCreatedAt            = "The creation date of the catalog."
	descCatalogDescription          = "The description of the catalog."
	descCatalogIsCached             = "Indicates whether the catalog is cached."
	descCatalogIsLocal              = "Indicates whether the catalog is local."
	descCatalogIsPublished          = "Indicates whether the catalog is published."
	descCatalogIsShared             = "Indicates whether the catalog is shared."
	descCatalogMediaItemList        = "The list of media items in the catalog."
	descCatalogNumberOfMedia        = "The number of media in the catalog."
	descCatalogOwnerName            = "The owner name of the catalog."
	descCatalogPreserveIdentityInfo = "Include BIOS UUIDs and MAC addresses in the downloaded OVF package. Keep in mind that preserving this identity information reduces the package's portability, so only include it when necessary."
	descMediaID                     = "The ID of the media."
	descMediaCreatedAt              = "The date and time when the media was created."
	descMediaDescription            = "The description of the media."
	descMediaIsISO                  = "`True` if the media is an ISO."
	descMediaIsPublished            = "`True` if the media is published."
	descMediaName                   = "The name of the media."
	descMediaOwnerName              = "The name of the owner of the media."
	descMediaSize                   = "The size of the media in bytes."
	descMediaStatus                 = "The status of the media."
	descMediaStorageProfile         = "The storage profile of the media."
)

type catalog interface {
	GetID() string
	GetName() string
	GetIDOrName() string
	GetCatalog() (*govcd.AdminCatalog, error)
}

type base struct {
	id   string
	name string
}

func (b base) GetID() string {
	return b.id
}

func (b base) GetName() string {
	return b.name
}

func (b base) GetIDOrName() string {
	if b.id != "" {
		return b.id
	}
	return b.name
}
