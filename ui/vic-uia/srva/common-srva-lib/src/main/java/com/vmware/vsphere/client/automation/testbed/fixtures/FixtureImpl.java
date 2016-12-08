/**
 * Copyright 2013 VMware, Inc.  All rights reserved. -- VMware Confidential
 */
package com.vmware.vsphere.client.automation.testbed.fixtures;

import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import com.vmware.client.automation.common.spec.ServiceSpec;
import com.vmware.vsphere.client.automation.common.spec.ManagedEntitySpec;
import com.vmware.vsphere.client.automation.srv.common.model.LibraryCommonConstants.LibraryType;
import com.vmware.vsphere.client.automation.srv.common.spec.ClusterSpec;
import com.vmware.vsphere.client.automation.srv.common.spec.DatacenterSpec;
import com.vmware.vsphere.client.automation.srv.common.spec.DatastoreSpec;
import com.vmware.vsphere.client.automation.srv.common.spec.HostSpec;
import com.vmware.vsphere.client.automation.srv.common.spec.LibrarySpec;
import com.vmware.vsphere.client.automation.srv.common.spec.SpecFactory;
import com.vmware.vsphere.client.automation.srv.common.spec.VcSpec;
import com.vmware.vsphere.client.automation.testbed.common.CommonTestBed;
import com.vmware.vsphere.client.automation.testbed.common.LocalTestBedImpl;
import com.vmware.vsphere.client.automation.testbed.model.FixtureEntities;
import com.vmware.vsphere.client.automation.testbed.model.Fixtures;
import com.vmware.vsphere.client.automation.testbed.spec.FixtureClusterSpec;
import com.vmware.vsphere.client.automation.testbed.spec.FixtureDatacenterSpec;
import com.vmware.vsphere.client.automation.testbed.spec.FixtureDatastoreSpec;
import com.vmware.vsphere.client.automation.testbed.spec.FixtureHostSpec;
import com.vmware.vsphere.client.automation.testbed.spec.FixtureLibrarySpec;
import com.vmware.vsphere.client.automation.testbed.spec.FixtureVcSpec;

/**
 * Class that represents the base common inventory setup
 * - Datacenter
 * - Cluster
 * - Host
 * - Datastore
 */
@Deprecated
public class FixtureImpl implements Fixture {

   private static final Logger _logger = LoggerFactory.getLogger(FixtureImpl.class);

   private final Fixtures _fixture;
   private final CommonTestBed _testbedUtil;

   /**
    * Factory method that creates a specified Fixture.
    * @param fixture - denotes which fixture to be created
    * @return Fixture instance that corresponds to the specified fixture
    */
   public static Fixture createFixture(Fixtures fixture, String propertyFile) {
      return new FixtureImpl(fixture, propertyFile);
   }

   private FixtureImpl(Fixtures fixture, String propertyFile) {
      _fixture = fixture;
      // TODO currently it is hardcoded that local setup is used;
      // add functionality that nimbus is supported - probably when
      // testbed is made a parameter of initSpec() method
      // PR 1129635
      _testbedUtil = new LocalTestBedImpl(propertyFile);
   }

   /**
    * {@inheritDoc}
    */

   @SuppressWarnings("unchecked")
   @Override
   public <T extends ManagedEntitySpec> T getFixtureResource(
         Class<T> specClass, FixtureEntities entity) {

      if (!_fixture.containsFixtureEntity(entity)) {
         throw new IllegalArgumentException(
               "The fixture doesn't contain the entity: " + entity);
      }

      if (specClass.getSimpleName().equals(VcSpec.class.getSimpleName())) {
         return (T) getVc(entity);
      } else if (specClass.getSimpleName().equals(HostSpec.class.getSimpleName())) {
         return (T) getHost(entity);
      } else if (specClass.getSimpleName().equals(ClusterSpec.class.getSimpleName())) {
         return (T) getCluster(entity);
      } else if (specClass.getSimpleName().equals(DatacenterSpec.class.getSimpleName())) {
         return (T) getDatacenter(entity);
      } else if (specClass.getSimpleName().equals(DatastoreSpec.class.getSimpleName())) {
         return (T) getDatastore(entity);
      } else if (specClass.getSimpleName().equals(LibrarySpec.class.getSimpleName())) {
         return (T) getContentLibrary(entity);
      }

      return null;
   }

   /**
    * Constructs FixtureVcSpec.
    * @param entity fixture entity that describes VC.
    * @return FixtureVcSpec.
    * @throws IllegalArgumentException if this entity doesn't correspond to a VC.
    */
   public FixtureVcSpec getVc(FixtureEntities vc) {
      switch (vc) {
      case VC:
         FixtureVcSpec vcSpec = SpecFactory.getSpec(FixtureVcSpec.class);
         vcSpec.name.set(_testbedUtil.getVcName());
         vcSpec.loginUsername.set(_testbedUtil.getVcUsername());
         vcSpec.loginPassword.set(_testbedUtil.getVcPassword());
         vcSpec.ssoLoginUsername.set(_testbedUtil.getVcSsoUsername());
         vcSpec.ssoLoginPassword.set(_testbedUtil.getVcSsoPassword());
         //vcSpec.thumbprint.set(_testbedUtil.getVcThumbprint());
         return vcSpec;
      default:
         throw new IllegalArgumentException("No such VC choice.");
      }
   }

   /**
    * Method that constructs a FixtureHostSpec with name and parent
    * @param host - fixture entity that describes host
    * @return - FixtureHostSpec with name and parent
    * @throws IllegalArgumentException if this entity doesn't correspond to a host
    */
   public FixtureHostSpec getHost(FixtureEntities host) {

      switch (host) {
      case NGC_COMMON_CLUSTERED_HOST:
         return SpecFactory.getSpec(FixtureHostSpec.class,
               _testbedUtil.getCommonHostName(),
               getCluster(FixtureEntities.NGC_COMMON_CLUSTER));
      default:
         throw new IllegalArgumentException("No such host choice.");
      }
   }

   /**
    * Method that constructs a FixtureClusterSpec with name and parent
    * @param host - fixture entity that describes cluster
    * @return - FixtureClusterSpec with name and parent
    * @throws IllegalArgumentException if this entity doesn't correspond to a cluster
    */
   public FixtureClusterSpec getCluster(FixtureEntities cluster) {

      switch (cluster) {
      case NGC_COMMON_CLUSTER:
         return SpecFactory.getSpec(FixtureClusterSpec.class,
               _testbedUtil.getCommonClusterName(),
               getDatacenter(FixtureEntities.NGC_COMMON_DATACENTER));
      default:
         throw new IllegalArgumentException("No such cluster choice.");
      }
   }

   /**
    * Method that constructs a FixtureDatacenterSpec with name and parent
    * @param host - fixture entity that describes datacenter
    * @return - FixtureDatacenterSpec with name and parent
    * @throws IllegalArgumentException if this entity doesn't correspond to a datacenter
    */
   public FixtureDatacenterSpec getDatacenter(FixtureEntities datacenter) {

      switch (datacenter) {
      case NGC_COMMON_DATACENTER:
         return SpecFactory.getSpec(FixtureDatacenterSpec.class,
               _testbedUtil.getCommonDatacenterName(), null);
      default:
         throw new IllegalArgumentException("No such datacenter choice.");
      }
   }

   /**
    * Method that constructs a FixtureDatastoreSpec with name and parent
    * @param host - fixture entity that describes a datastore
    * @return - FixtureDatastoreSpec with name and parent
    * @throws IllegalArgumentException if this entity doesn't correspond to a datastore
    */
   public FixtureDatastoreSpec getDatastore(FixtureEntities datastore) {

      switch (datastore) {
      case NGC_COMMON_DATASTORE:
         return SpecFactory.getSpec(FixtureDatastoreSpec.class,
               _testbedUtil.getCommonDatastoreName(),
               getDatacenter(FixtureEntities.NGC_COMMON_DATACENTER));
      default:
         throw new IllegalArgumentException("No such datastore choice.");
      }
   }

   /**
    * Method that constructs a FixtureLibrarySpec with name and parent
    * @param host - fixture entity that describes a content library
    * @return - FixtureLibrarySpec with name and parent
    * @throws Exception if this entity doesn't correspond to a content library
    */
   public FixtureLibrarySpec getContentLibrary(FixtureEntities contentLibrary) {
      // Will be linked to the LibrarySpec.
      ServiceSpec serviceSpec = toServiceSpec(getVc(FixtureEntities.VC));

      switch (contentLibrary) {
      case CONTENT_LIBRARY_LOCAL:
         FixtureLibrarySpec localLibrarySpec =
         SpecFactory.getSpec(
               FixtureLibrarySpec.class,
               _testbedUtil.getContentLibraryName(),
               null);
         localLibrarySpec.libraryType.set(LibraryType.LOCAL);
         localLibrarySpec.datastore.set(
               getDatastore(FixtureEntities.NGC_COMMON_DATASTORE));

         localLibrarySpec.links.add(serviceSpec);
         return localLibrarySpec;
      case CONTENT_LIBRARY_PUBLISHED:
         FixtureLibrarySpec localPublishedLibrarySpec =
         SpecFactory.getSpec(
               FixtureLibrarySpec.class,
               _testbedUtil.getContentLibraryName(),
               null);
         localPublishedLibrarySpec.libraryType.set(LibraryType.LOCAL);
         localPublishedLibrarySpec.isPublished.set(true);
         localPublishedLibrarySpec.datastore.set(
               getDatastore(FixtureEntities.NGC_COMMON_DATASTORE));

         localPublishedLibrarySpec.links.add(serviceSpec);
         return localPublishedLibrarySpec;
      case CONTENT_LIBRARY_SUBSCRIBED:
         FixtureLibrarySpec subscribedLibrarySpec =
         SpecFactory.getSpec(
               FixtureLibrarySpec.class,
               _testbedUtil.getContentLibraryName(),
               null);
         subscribedLibrarySpec.libraryType.set(LibraryType.SUBSCRIBED);
         subscribedLibrarySpec.datastore.set(
               getDatastore(FixtureEntities.NGC_COMMON_DATASTORE));

         subscribedLibrarySpec.links.add(serviceSpec);
         return subscribedLibrarySpec;
      default:
         throw new IllegalArgumentException(
               "No such content library choice.");
      }
   }

   public static ServiceSpec toServiceSpec(VcSpec vcSpec) {
      ServiceSpec serviceSpec = new ServiceSpec();
      String errMsgFormat = "VcSpec lacks '%s' assigned!";

      if (!vcSpec.name.isAssigned()) {
         throw new IllegalArgumentException(String.format(errMsgFormat, "name"));
      }
      if (!vcSpec.loginUsername.isAssigned()) {
         throw new IllegalArgumentException(String.format(errMsgFormat, "loginUsername"));
      }
      if (!vcSpec.loginPassword.isAssigned()) {
         throw new IllegalArgumentException(String.format(errMsgFormat, "loginPassword"));
      }

      serviceSpec.endpoint.set(vcSpec.name.get());

      if (vcSpec.ssoLoginUsername.isAssigned()) {
         serviceSpec.username.set(vcSpec.ssoLoginUsername.get());
      } else {
         _logger.warn(String.format(errMsgFormat, "ssoLoginUsername"));
      }
      if (vcSpec.ssoLoginPassword.isAssigned()) {
         serviceSpec.password.set(vcSpec.ssoLoginPassword.get());
      } else {
         _logger.warn(String.format(errMsgFormat, "ssoLoginPassword"));
      }

      return serviceSpec;
   }
}
