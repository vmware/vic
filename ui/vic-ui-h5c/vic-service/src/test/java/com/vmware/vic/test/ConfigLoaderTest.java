/*

Copyright 2017 VMware, Inc. All Rights Reserved.

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
package com.vmware.vic.test;

import static org.junit.Assert.*;

import java.io.FileNotFoundException;
import java.io.IOException;

import org.junit.Test;

import com.vmware.vic.utils.ConfigLoader;

@SuppressWarnings("unused")
public class ConfigLoaderTest {
	@Test
	public void testConfigLoaderLoads() {
		try {
			ConfigLoader configLoader = new ConfigLoader("configs.properties");
			assertNotNull(configLoader);
		} catch (IOException ex) {
			ex.printStackTrace();
		}
	}

	@Test
	public void testWrongConfigsFile() {
		try {
			ConfigLoader configLoader = new ConfigLoader("meh");
		} catch (FileNotFoundException ex) {
			assertNotNull(ex.getMessage());
		} catch (IOException e) {
			assertNotNull(e);
		}
	}

	@Test
	public void testGetPropertyUiVersion() throws Exception {
		ConfigLoader configLoader = new ConfigLoader("configs.properties");
		assertTrue(configLoader.getProp("uiVersion").length() > 0);
	}
}
