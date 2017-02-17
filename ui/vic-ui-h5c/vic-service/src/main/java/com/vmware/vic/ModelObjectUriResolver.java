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
package com.vmware.vic;

import java.net.URI;
import java.net.URISyntaxException;

import com.vmware.vise.data.uri.ResourceTypeResolver;

public class ModelObjectUriResolver implements ResourceTypeResolver {
	private static final String SCHEME = "urn";
	private static final String NAMESPACE = "vic";
	private static final String UID_PREFIX = SCHEME + ":" + NAMESPACE;

	private static final String TYPE_DELIMITER = ":";
	private static final String FRAGMENT_SEPARATOR = "/";

	/*
	 * (non-Javadoc)
	 * @see com.vmware.vise.data.uri.ResourceTypeResolver#getResourceType(java.net.URI)
	 */
	@Override
	public String getResourceType(URI uri) {
		if (!isValid(uri)) {
			throwIllegalURIException(uri);
		}

		return parseUri(uri, true);
	}

	/*
	 * (non-Javadoc)
	 * @see com.vmware.vise.data.uri.ResourceTypeResolver#getServerGuid(java.net.URI)
	 */
	@Override
	public String getServerGuid(URI uri) {
		String id = parseUri(uri, false);
		int fragmentSeperatorIndex = id.indexOf(FRAGMENT_SEPARATOR);
		if (fragmentSeperatorIndex <= 0) {
			throwIllegalURIException(uri);
		}
		return id.substring(0, fragmentSeperatorIndex);
	}

	/**
	 * Parse URI object according to parseType given
	 * @param uri
	 * @param parseType
	 * @return resourceType if parseType is true
               resourceId if parseType is false
	 */
	private String parseUri(URI uri, boolean parseType) {
		// get substring after urn:
		String ssPart = uri.getSchemeSpecificPart();
		int typeIndex = ssPart.indexOf(TYPE_DELIMITER);
		ssPart = ssPart.substring(typeIndex + 1);
		int resourceIndex = ssPart.lastIndexOf(TYPE_DELIMITER);

		if (resourceIndex == -1) {
			throw new IllegalArgumentException(
					"Invalid URI. Missing type delimiter " + toString(uri));
		}

		String result;
		if (parseType) {
			result = ssPart.substring(0, resourceIndex);
		} else {
			result = ssPart.substring(resourceIndex + 1);
		}

		return result;
	}

	/**
	 * Generate a URI instance with resource type and id
	 * @param type
	 * @param id
	 * @return URI
	 */
	public URI createUri(String type, String id) {
		if (type == null || type.length() < 1) {
			throw new IllegalArgumentException("type must be non-null");
		}

		if (id == null || id.length() < 1) {
			throw new IllegalArgumentException("id must be non-null");
		}

		URI uri = null;
		try {
			String schemeSpecificPart =
					NAMESPACE + TYPE_DELIMITER + type + TYPE_DELIMITER + id;
			uri = new URI(SCHEME, schemeSpecificPart, null);
		} catch (URISyntaxException e) {
			throw new IllegalArgumentException(e);
		}
		return uri;
	}

	/**
	 * Return resourceId
	 * @param uri
	 * @return resource id of URI (e.g. server1/vic-root
			   from URI urn:vic:vic:root:server1/vic-root)
	 */
	public String getId(URI uri) {
		if (!isValid(uri)) {
			throwIllegalURIException(uri);
		}
		return parseUri(uri, false);
	}

	/**
	 * Return objectId of resourceId
	 * @param uri
	 * @return object id of resourceId (e.g. vic-root from
	 *         server1/vic-root)
	 */
	public String getObjectId(URI uri) {
		String id = parseUri(uri, false);
		int fragmentSeperatorIndex = id.indexOf(FRAGMENT_SEPARATOR);
		if (fragmentSeperatorIndex <= 0) {
			throwIllegalURIException(uri);
		}
		return id.substring(fragmentSeperatorIndex + 1);
	}

	/**
	 * Return the string representation of the URI
	 * @param uri
	 * @return uri.toString()
	 */
	public String getUid(URI uri) {
		if (!isValid(uri)) {
			throwIllegalURIException(uri);
		}
		return uri.toString();
	}

	/**
	 * Create and return the resourceId based on serverGuid and objectId
	 * @param serverGuid
	 * @param objectId
	 * @return serverGuid appended by FRAGMENT_SEPARATOR and objectId such that
	 * 		   the object can be used uniquely identified within any server
	 */
	public String createResourceId(String serverGuid, String objectId) {
		return serverGuid + FRAGMENT_SEPARATOR + objectId;
	}

	private boolean isValid(URI uri) {
		return (uri != null) && uri.toString().startsWith(UID_PREFIX);
	}

	private void throwIllegalURIException(URI uri) {
		throw new IllegalArgumentException("URI " +
				toString(uri) + " is invalid for this resolver");
	}

	private String toString(URI uri) {
		if (uri == null) {
			return null;
		}
		return uri.toString();
	}
}
