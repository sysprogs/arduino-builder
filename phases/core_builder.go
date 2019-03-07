/*
 * This file is part of Arduino Builder.
 *
 * Arduino Builder is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 2 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program; if not, write to the Free Software
 * Foundation, Inc., 51 Franklin St, Fifth Floor, Boston, MA  02110-1301  USA
 *
 * As a special exception, you may use this file as part of a free software
 * library without restriction.  Specifically, if other files instantiate
 * templates or use macros or inline functions from this file, or you compile
 * this file and link it with other files to produce an executable, this
 * file does not by itself cause the resulting executable to be covered by
 * the GNU General Public License.  This exception does not however
 * invalidate any other reasons why the executable file might be covered by
 * the GNU General Public License.
 *
 * Copyright 2015 Arduino LLC (http://www.arduino.cc/)
 */

package phases

import (
	"os"
	"path/filepath"

	"github.com/arduino/arduino-builder/builder_utils"
	"github.com/arduino/arduino-builder/constants"
	"github.com/arduino/arduino-builder/i18n"
	"github.com/arduino/arduino-builder/types"
	"github.com/arduino/arduino-builder/utils"
	"github.com/arduino/go-properties-map"
)

type CoreBuilder struct{}

func (s *CoreBuilder) Run(ctx *types.Context) error {
	coreBuildPath := ctx.CoreBuildPath
	coreBuildCachePath := ctx.CoreBuildCachePath
	var buildProperties = ctx.BuildProperties

	err := utils.EnsureFolderExists(coreBuildPath)
	if err != nil {
		return i18n.WrapError(err)
	}

	if coreBuildCachePath != "" {
		err := utils.EnsureFolderExists(coreBuildCachePath)
		if err != nil {
			return i18n.WrapError(err)
		}
	}

	var coreModel *types.CodeModelLibrary
	if ctx.CodeModelBuilder != nil {
		coreModel = new(types.CodeModelLibrary)
		ctx.CodeModelBuilder.Core = coreModel
	}

	if ctx.UnoptimizeCore {
		buildProperties = builder_utils.RemoveOptimizationFromBuildProperties(buildProperties)
	}
	
	buildProperties = builder_utils.ExpandSysprogsExtensionProperties(buildProperties)

	archiveFile, objectFiles, err := compileCore(ctx, coreBuildPath, coreBuildCachePath, buildProperties, coreModel)
	if err != nil {
		return i18n.WrapError(err)
	}

	ctx.CoreArchiveFilePath = archiveFile
	ctx.CoreObjectsFiles = objectFiles

	return nil
}

func compileCore(ctx *types.Context, buildPath string, buildCachePath string, buildProperties properties.Map, coreModel *types.CodeModelLibrary) (string, []string, error) {
	logger := ctx.GetLogger()
	coreFolder := buildProperties[constants.BUILD_PROPERTIES_BUILD_CORE_PATH]
	variantFolder := buildProperties[constants.BUILD_PROPERTIES_BUILD_VARIANT_PATH]

	targetCoreFolder := buildProperties[constants.BUILD_PROPERTIES_RUNTIME_PLATFORM_PATH]

	if coreModel != nil {
		coreModel.SourceDirectory = coreFolder
		coreModel.Name = buildProperties[constants.LIBRARY_NAME]
	}

	includes := []string{}
	includes = append(includes, coreFolder)
	if variantFolder != constants.EMPTY_STRING {
		includes = append(includes, variantFolder)
	}
	includes = utils.Map(includes, utils.WrapWithHyphenI)

	var err error

	variantObjectFiles := []string{}
	if variantFolder != constants.EMPTY_STRING {
		variantObjectFiles, err = builder_utils.CompileFiles(ctx, variantObjectFiles, variantFolder, true, buildPath, buildProperties, includes, coreModel)
		if err != nil {
			return "", nil, i18n.WrapError(err)
		}
	}

	// Recreate the archive if ANY of the core files (including platform.txt) has changed
	realCoreFolder := utils.GetParentFolder(coreFolder, 2)

	var targetArchivedCore string
	if buildCachePath != "" {
		archivedCoreName := builder_utils.GetCachedCoreArchiveFileName(buildProperties[constants.BUILD_PROPERTIES_FQBN], realCoreFolder)
		targetArchivedCore = filepath.Join(buildCachePath, archivedCoreName)
		canUseArchivedCore := !builder_utils.CoreOrReferencedCoreHasChanged(realCoreFolder, targetCoreFolder, targetArchivedCore)

		if canUseArchivedCore {
			// use archived core
			if ctx.Verbose {
				logger.Println(constants.LOG_LEVEL_INFO, "Using precompiled core: {0}", targetArchivedCore)
			}
			return targetArchivedCore, variantObjectFiles, nil
		}
	}

	coreObjectFiles, err := builder_utils.CompileFiles(ctx, []string{}, coreFolder, true, buildPath, buildProperties, includes, coreModel)
	if err != nil {
		return "", nil, i18n.WrapError(err)
	}

	archiveFile, err := builder_utils.ArchiveCompiledFiles(ctx, buildPath, "core.a", coreObjectFiles, buildProperties, coreModel)
	if err != nil {
		return "", nil, i18n.WrapError(err)
	}

	// archive core.a
	if targetArchivedCore != "" && coreModel != nil {
		err := builder_utils.CopyFile(archiveFile, targetArchivedCore)
		if ctx.Verbose {
			if err == nil {
			logger.Println(constants.LOG_LEVEL_INFO, constants.MSG_ARCHIVING_CORE_CACHE, targetArchivedCore)
			} else if os.IsNotExist(err) {
				logger.Println(constants.LOG_LEVEL_INFO, constants.MSG_CORE_CACHE_UNAVAILABLE, ctx.ActualPlatform.PlatformId)
			} else {
				logger.Println(constants.LOG_LEVEL_INFO, constants.MSG_ERROR_ARCHIVING_CORE_CACHE, targetArchivedCore, err)
		}
	}
	}

	return archiveFile, variantObjectFiles, nil
}
