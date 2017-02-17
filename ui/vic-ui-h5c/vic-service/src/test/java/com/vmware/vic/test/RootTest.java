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

import java.net.URI;

import org.junit.Before;
import org.junit.Test;

import com.vmware.vic.ModelObjectUriResolver;
import com.vmware.vic.model.ModelObject;
import com.vmware.vic.model.Root;
import com.vmware.vic.model.RootInfo;

public class RootTest {
	private Root _root;

	@Before
	public void setModelObject() {
		_root = new Root(
				new RootInfo(new String[]{"1.0"}), 0, 0);
	}

	@Test
	public void testGetType() {
		assertEquals("vic:Root", _root.getType());
	}

	@Test
	public void testGetId() {
		assertEquals("vic/vic-root", _root.getId());
	}

	@Test
	public void testGetUri() {
		ModelObjectUriResolver uriResolver = new ModelObjectUriResolver();
		URI actual = _root.getUri(uriResolver);
		assertEquals("urn:vic:vic:Root:vic/vic-root", actual.toString());
	}

	@Test
	public void testGetProperty() {
		assertFalse(_root.getProperty("uiVersion").equals(ModelObject.UNSUPPORTED_PROPERTY));
		assertFalse(_root.getProperty("vchVmsLen").equals(ModelObject.UNSUPPORTED_PROPERTY));
		assertFalse(_root.getProperty("containerVmsLen").equals(ModelObject.UNSUPPORTED_PROPERTY));
		assertTrue(_root.getProperty("iDontExist").equals(ModelObject.UNSUPPORTED_PROPERTY));
	}

	@Test
	public void testToString() {
		assertTrue(_root.toString().contentEquals("uiVersion: 1.0, vchVms.length: 0, containerVms.length: 0"));
	}
}
