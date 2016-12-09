/* Copyright 2013 VMware, Inc. All rights reserved. -- VMware Confidential */
package com.vmware.vsphere.client.automation.srv.common.spec;

import com.vmware.hsua.common.datamodel.BasePBox.DataProperty;
import com.vmware.vsphere.client.automation.common.spec.ManagedEntitySpec;
import com.vmware.vsphere.client.automation.srv.common.model.LibraryCommonConstants.DownloadContentType;
import com.vmware.vsphere.client.automation.srv.common.model.LibraryCommonConstants.LibraryType;

/**
 * Container class for library properties.
 *
 * TODO: The proper way is to also remove the public libraryType property,
 * make the LibrarySpec abstract and introduce 3 descendant classes:
 * LocalLibrarySpec, PublishedLibrarySpec and SubscribedLibrarySpec.
 */
public class LibrarySpec extends ManagedEntitySpec {
   /**
    * Description of the object.
    */
   public DataProperty<String> description;

   /**
    * Type of the library. The library can be local or subscribed.
    */
   public DataProperty<LibraryType> libraryType;

   /**
    * Select if Streaming Optimized is required.
    */
   public DataProperty<Boolean> isStreamingOptimized;

   /**
    * The username used for published libraries.
    */
   public DataProperty<String> username;

   /**
    * The password used for published libraries.
    */
   public DataProperty<String> password;

   /**
    * The password used for subscribed libraries.
    */
   public DataProperty<String> subscriptionPassword;

   /**
    * Select if authentication is required.
    */
   public DataProperty<Boolean> authenticationRequired;

   /**
    * Whether content is downloaded immediately or on demand.
    */
   public DataProperty<DownloadContentType> downloadContentType;

   /**
    * Published status of the library. A local library can be Published or not.
    */
   public DataProperty<Boolean> isPublished;

   /**
    * The URL of the library to which to subscribe.
    */
   public DataProperty<String> subscriptionUrl;

   /**
    * The sslthumbprint of the VC where is situated the library to which to subscribe
    */
   public DataProperty<String> sslThumbprint;

   /**
    * Specification of the datastore that will be used for
    * storage backing of the content library
    */
   public DataProperty<DatastoreSpec> datastore;

   /**
    * Specification of the file system directory path that will be used for
    * storage backing of the content library
    */
   public DataProperty<String> directoryPath;
}
