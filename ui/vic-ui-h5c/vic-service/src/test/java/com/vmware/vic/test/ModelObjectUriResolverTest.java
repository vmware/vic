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
import java.net.URISyntaxException;

import org.junit.Before;
import org.junit.Test;

import com.vmware.vic.ModelObjectUriResolver;

public class ModelObjectUriResolverTest {
	private ModelObjectUriResolver moUriResolver;

	@Before
	public void setModelObjectUriResolver() {
		moUriResolver = new ModelObjectUriResolver();
	}

	@Test
	public void testGetResourceType() throws URISyntaxException {
		URI uri = new URI("urn", "vic:vic:Root:server1/rootObject", null);
		String resourceType = moUriResolver.getResourceType(uri);
		assertEquals("vic:Root", resourceType);
	}

	@Test
	public void testGetServerGuid() throws URISyntaxException {
		URI uri = new URI("urn", "vic:vic:Root:server1/rootObject", null);
		String serverGuid = moUriResolver.getServerGuid(uri);
		assertEquals("server1", serverGuid);
	}

	@Test
	public void testGetObjectId() throws URISyntaxException {
		URI uri = new URI("urn", "vic:vic:VirtualContainerHostVm:vic/ALL", null);
		String objectId = moUriResolver.getObjectId(uri);
		assertEquals("ALL", objectId);
	}

	@Test
	public void testCreateUri() {
		URI uri = moUriResolver.createUri("vic:Root", "server1/rootObject");
		assertEquals("urn:vic:vic:Root:server1/rootObject", uri.toString());
	}
}
