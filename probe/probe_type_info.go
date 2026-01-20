/*
Copyright 2024 4rcadia

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package probe

import "time"

// MinecraftServerDotProperties is not universal across game versions. Therefore,
// it is just a map[string]string. Just remember to check the game version before
// calling any newly added property.
type MinecraftServerDotProperties map[string]string

// VersionDotJson (version.json) do not exist across all game versions.
// TODO: Backtrack the last version with this file
// TODO: Implement an alternative version check method (checksum)
type VersionDotJson struct {
	Id              string `json:"id"`
	Name            string `json:"name"`
	WorldVersion    int    `json:"world_version"`
	SeriesId        string `json:"series_id"`
	ProtocolVersion int    `json:"protocol_version"`
	PackVersion     struct {
		Resource int `json:"resource"`
		Data     int `json:"data"`
	} `json:"pack_version"`
	BuildTime     time.Time `json:"build_time"`
	JavaComponent string    `json:"java_component"`
	JavaVersion   int       `json:"java_version"`
	Stable        bool      `json:"stable"`
	UseEditor     bool      `json:"use_editor"`
}

type McdrConfigDotYml struct {
	Language         string `yaml:"language"`
	WorkingDirectory string `yaml:"working_directory"`
	StartCommand     string `yaml:"start_command"`
	Handler          string `yaml:"handler"`
	Encoding         string `yaml:"encoding"`
	Decoding         string `yaml:"decoding"`
	Rcon             struct {
		Enable   bool   `yaml:"enable"`
		Address  string `yaml:"address"`
		Port     int    `yaml:"port"`
		Password string `yaml:"password"`
	} `yaml:"rcon"`
	PluginDirectories         []string    `yaml:"plugin_directories"`
	CatalogueMetaCacheTtl     int         `yaml:"catalogue_meta_cache_ttl"`
	CatalogueMetaFetchTimeout int         `yaml:"catalogue_meta_fetch_timeout"`
	CatalogueMetaUrl          interface{} `yaml:"catalogue_meta_url"`
	PluginDownloadUrl         interface{} `yaml:"plugin_download_url"`
	PluginDownloadTimeout     int         `yaml:"plugin_download_timeout"`
	PluginPipInstallExtraArgs interface{} `yaml:"plugin_pip_install_extra_args"`
	CheckUpdate               bool        `yaml:"check_update"`
	AdvancedConsole           bool        `yaml:"advanced_console"`
	HttpProxy                 interface{} `yaml:"http_proxy"`
	HttpsProxy                interface{} `yaml:"https_proxy"`
	Telemetry                 bool        `yaml:"telemetry"`
	DisableConsoleThread      bool        `yaml:"disable_console_thread"`
	DisableConsoleColor       bool        `yaml:"disable_console_color"`
	CustomHandlers            interface{} `yaml:"custom_handlers"`
	CustomInfoReactors        interface{} `yaml:"custom_info_reactors"`
	WatchdogThreshold         int         `yaml:"watchdog_threshold"`
	HandlerDetection          bool        `yaml:"handler_detection"`
	Debug                     struct {
		All          bool `yaml:"all"`
		Mcdr         bool `yaml:"mcdr"`
		Handler      bool `yaml:"handler"`
		Reactor      bool `yaml:"reactor"`
		Plugin       bool `yaml:"plugin"`
		Permission   bool `yaml:"permission"`
		Command      bool `yaml:"command"`
		TaskExecutor bool `yaml:"task_executor"`
		Telemetry    bool `yaml:"telemetry"`
	} `yaml:"debug"`
}
