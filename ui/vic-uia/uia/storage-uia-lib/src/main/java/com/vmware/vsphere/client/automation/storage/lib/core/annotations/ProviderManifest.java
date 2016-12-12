package com.vmware.vsphere.client.automation.storage.lib.core.annotations;

import java.lang.annotation.ElementType;
import java.lang.annotation.Retention;
import java.lang.annotation.RetentionPolicy;
import java.lang.annotation.Target;

import com.vmware.vsphere.client.automation.storage.lib.core.tests.BaseTest;

/**
 * Annotation for adding multiple @{@link UsesProvider} annotation on class.
 * <p>
 * see {@link UsesProvider}
 * </p>
 * <p>
 * If a class is deriving from a {@link BaseTest} and has this annotation all of
 * the {@link UsesProvider} will be added to the tests workflow specs
 * </p>
 */
@Retention(RetentionPolicy.RUNTIME)
@Target({ ElementType.TYPE })
public @interface ProviderManifest {
   UsesProvider[] value();
}
