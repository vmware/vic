package com.vmware.vicui.views {

	import com.vmware.core.model.IResourceReference;
	import com.vmware.data.query.DataUpdateSpec;
	import com.vmware.data.query.events.DataByModelRequest;
	import com.vmware.data.query.events.DataRequestInfo;
	import com.vmware.ui.IContextObjectHolder;
	import com.vmware.vicui.constants.AppConstants;
	import com.vmware.vicui.model.ContainerInfo;
	
	import flash.events.EventDispatcher;
	
	import mx.logging.ILogger;
	import mx.logging.Log;
	
	[Event(name="{com.vmware.data.query.events.DataByModelRequest.REQUEST_ID}",
   type="com.vmware.data.query.events.DataByModelRequest")]

	/**
	 * The mediator for ContainerPortletView
	 */
	public class ContainerPortletMediator extends EventDispatcher implements IContextObjectHolder
	{
	   private var _contextObject:IResourceReference;
	   private var _view:ContainerPortletView;
	
	   private static var _logger:ILogger = Log.getLogger('VicuiMediator');
	
	   [View]
	   /** The view associated with this mediator. */
	   public function set view(value:ContainerPortletView):void {
		   _view = value;
	   }
	   
	   /**
		* Returns the view.
		*/
	   public function get view():ContainerPortletView {
		   return _view;
	   }
	
	
	   [Bindable]
	   /** Returns the inventory object handled in this view (IContextObjectHolder interface) */
	   public function get contextObject():Object {
	      return _contextObject;
	   }
	
	   /** Called by the framework with the current inventory object or null */
	   public function set contextObject(value:Object):void {
	      _contextObject = IResourceReference(value);
	
	      if (_contextObject == null) {
	         // A null contextObject means that the view is being cleared
	         clearData();
	         return;
	      }
	
	      // Once contextObject is set the view can be initialized with the object data.
	      requestData();
	   }
	
	   private function requestData():void {
	   	   // Default data request option allowing implicit updates of the view
	   	   var requestInfo:DataRequestInfo = new DataRequestInfo(DataUpdateSpec.newImplicitInstance());

		   // Dispatch an event to fetch the _contextObject data from the server along the specified model.
		   dispatchEvent(DataByModelRequest.newInstance(_contextObject, ContainerInfo, requestInfo));
	   }
	   
	   [ResponseHandler(name="{com.vmware.data.query.events.DataByModelRequest.RESPONSE_ID}")]
	   public function onData(event:DataByModelRequest, result:ContainerInfo):void {
		   _logger.info("Container summary data retrieved.");
		   
		   if (_view != null) {

			   //set default placeholder data
			   _view.isContainer = new Boolean(false);
			   _view.hasPortmappingInfo = new Boolean(false);
			   _view.containerName.text = new String(AppConstants.PLACEHOLDER_VAL);
			   _view.imageName.text = new String(AppConstants.PLACEHOLDER_VAL);
			   _view.portmappingInfo.text = new String(AppConstants.PLACEHOLDER_VAL);

			   if (result != null) {
				   
				   var config:Array = new Array();
				   
				   //extraConfig data from vm config
				   config = result.extraConfig;
				   
				   if (config != null) {

					   var keyName:String = new String();
					   var keyVal:String = new String();
					   var key:String = new String();
					    
					   for (key in config) {
						   
					    	keyName = config[key].key.value as String;
							keyVal = config[key].value as String;
							
							//determine if this is container vm
							if (keyName == AppConstants.VM_CONTAINER_NAME_PATH) {
							    _view.isContainer = true;
							    _view.containerName.text = keyVal;
							    continue;
							}
							
							//get container image name
							if (keyName == AppConstants.VM_CONTAINER_IMAGE_PATH) {
							    _view.imageName.text = keyVal;
							   	continue;
							}

							//get container port mapping information
							if (keyName == AppConstants.VM_CONTAINER_PORTMAPPING) {
								_view.hasPortmappingInfo = true;
								_view.portmappingInfo.text = keyVal.replace(/\|/g, ", ");
								continue;
							}
					   }
				   }
			  }
		   }
	   }
	   
	   private function clearData() : void {
	   	   if (_view != null) {
		      // clear the UI data
			   _view.isContainer = false;
			   _view.hasPortmappingInfo = false;
			   _view.containerName.text = new String(AppConstants.PLACEHOLDER_VAL);
			   _view.imageName.text = new String(AppConstants.PLACEHOLDER_VAL);
			   _view.portmappingInfo.text = new String(AppConstants.PLACEHOLDER_VAL);
		   }
	   }
	}
}