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
package com.vmware.vic.model;

import java.net.URI;

import com.vmware.vic.ModelObjectUriResolver;

public abstract class ModelObject {
	public static final String NAMESPACE = "vic:";
	public static final Object UNSUPPORTED_PROPERTY = new Object();

	private String _id;
	private URI _uri;
	private String _type;

	/**
	 * @return object id
	 */
	public String getId() {
		return _id;
	}

	/**
	 * @param value
	 */
	protected void setId(String value) {
		_id = value;
	}

	/**
	 * @return object type
	 */
	public String getType() {
		if (_type == null) {
			_type = NAMESPACE + this.getClass().getSimpleName();
		}
		return _type;
	}

	/**
	 * @param resolver
	 * @return URI for this object
	 */
	public URI getUri(ModelObjectUriResolver resolver) {
		if (_uri == null) {
			_uri = resolver.createUri(getType(), _id);
		}
		return _uri;
	}

	public abstract Object getProperty(String property);
}
