package com.vmware.vsphere.client.automation.storage.lib.core.annotations;

import java.lang.annotation.ElementType;
import java.lang.annotation.Retention;
import java.lang.annotation.RetentionPolicy;
import java.lang.annotation.Target;

import com.vmware.client.automation.workflow.provider.ProviderWorkflow;
import com.vmware.vsphere.client.automation.storage.lib.core.tests.BaseTest;
import com.vmware.vsphere.client.automation.storage.lib.core.tests.ProviderPublication.ProviderPublicationUsage;

/**
 * Annotation for field of class to declare usage of Provider
 *
 * <p>
 * When the {@link UsesProvider} is used on of class deriving from
 * {@link BaseTest} field the corresponding provider entity will be pushed to
 * tests workflow specs and the field will be assigned to the matching spec
 * </p>
 *
 * <p>
 * When the {@link UsesProvider} is used as part {@link ProviderManifest}
 * annotation of class level for class deriving from {@link BaseTest} the
 * matching provider entities will be pushed to the tests workflow specs
 * </p>
 *
 */
@Retention(RetentionPolicy.RUNTIME)
@Target({ ElementType.FIELD })
public @interface UsesProvider {

   /**
    * The String id of the Provider
    *
    * @return
    */
   String id();

   /**
    * The class for the provider
    *
    * @return
    */
   Class<? extends ProviderWorkflow> clazz();

   /**
    * The String id with which the provider is publishing the entity
    *
    * @return
    */
   String entity();

   /**
    * The usage of the published entity.
    *
    * @return
    */
   ProviderPublicationUsage usage() default ProviderPublicationUsage.EXCLUSIVE;
}
