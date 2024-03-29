/*
Copyright 2020 WILDCARD

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
Created on 20/03/2021
*/

package config

type Hook struct {
	URL   string `json:"url"  yaml:"url"`
	Scope string `json:"scope" yaml:"scope"`
}

type Config struct {
	PreScriptFolder  string `json:"pre_script_folder" yaml:"pre_script_folder"`
	MainScriptFolder string `json:"main_script_folder" yaml:"main_script_folder"`
	PostScriptFolder string `json:"post_script_folder" yaml:"post_script_folder"`
	Hooks            []Hook `json:"hooks" yaml:"hooks"`
}

var (
	config     = new(Config)
	preScript  []string
	mainScript []string
	postScript []string
)
