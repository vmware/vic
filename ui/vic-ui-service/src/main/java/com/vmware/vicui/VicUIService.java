package com.vmware.vicui;

/**
 * Interface used to test a java service call from the Flex UI.
 *
 * It must be declared as osgi:service with the same name in
 * main/resources/META-INF/spring/bundle-context-osgi.xml
 */

import com.vmware.vise.data.query.PropertyProviderAdapter;

public interface VicUIService extends PropertyProviderAdapter {
   
}
