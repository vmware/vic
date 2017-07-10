/**
 * Copyright (c) 2016 VMware, Inc. All rights reserved.
 */
package com.vmware.vicui.model {
   import com.vmware.core.model.DataObject;

   [Bindable]
   //Declares the type of object associated with this model.
   //Types must be qualified with a namespace.
   [Model(type="VirtualMachine")]
   public class ContainerInfo extends DataObject {

      // Maps a model property to the class field.
      // Note: Property and field names don't need to match.
      // Also, properties within your own type don't require a namespace.
      /**
      * Container extra config data
	  * 
      */

	   [Model(property="config.extraConfig")]
	   public var extraConfig:Array;
   }
}
