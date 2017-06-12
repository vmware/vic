package com.vmware.vicui.views {

	import com.vmware.core.model.IResourceReference;
	import com.vmware.data.query.DataUpdateSpec;
	import com.vmware.data.query.events.DataByModelRequest;
	import com.vmware.data.query.events.DataRequestInfo;
	import com.vmware.ui.IContextObjectHolder;
	import com.vmware.vicui.constants.AppConstants;
	import com.vmware.vicui.model.VchInfo;

	import flash.events.EventDispatcher;
	import flash.utils.ByteArray;

	import mx.logging.ILogger;
	import mx.logging.Log;
	import mx.utils.Base64Decoder;

	[Event(name="{com.vmware.data.query.events.DataByModelRequest.REQUEST_ID}",
   type="com.vmware.data.query.events.DataByModelRequest")]

	/**
	 * The mediator for ContainerPortletView
	 */
	public class VchPortletMediator extends EventDispatcher implements IContextObjectHolder
	{
	   private var _contextObject:IResourceReference;
	   private var _view:VchPortletView;

	   private static var _logger:ILogger = Log.getLogger('VchMediator');

	   [View]
	   /** The view associated with this mediator. */
	   public function set view(value:VchPortletView):void {
		   _view = value;
	   }

	   /**
		* Returns the view.
		*/
	   public function get view():VchPortletView {
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
		   dispatchEvent(DataByModelRequest.newInstance(_contextObject, VchInfo, requestInfo));
	   }

	   [ResponseHandler(name="{com.vmware.data.query.events.DataByModelRequest.RESPONSE_ID}")]
	   public function onData(event:DataByModelRequest, result:VchInfo):void {
		   _logger.info("Vch summary data retrieved.");

		   if (_view != null) {

			   //set default placeholder data
			   _view.isVch = new Boolean(false);
			   _view.dockerApiEndpoint.text = new String(AppConstants.PLACEHOLDER_VAL);
			   _view.dockerLog.label = new String(AppConstants.PLACEHOLDER_VAL);

			   if (result != null) {

				   var config:Array = new Array();

				   //extraConfig data from vm config
				   config = result.extraConfig;

				   if (config != null) {

					   var keyName:String = new String();
					   var keyVal:String = new String();
					   var key:String = new String();
					   var isUsingTls:Boolean = true;

					   for (key in config ) {

						   keyName = config[key].key.value as String;
						   keyVal = config[key].value as String;

						   //determine if its a vch
						   if (keyName == AppConstants.VCH_NAME_PATH) {
							   _view.isVch = true;
							   continue;
						   }

						   //get container ip and decode to correct format
						   if (keyName == AppConstants.VCH_CLIENT_IP_PATH ) {
							   var base64Decoder:Base64Decoder = new Base64Decoder();
							   base64Decoder.decode(keyVal);

							   var bytes:ByteArray = new ByteArray();
							   var ip_ipv4:String = new String();

							   bytes = base64Decoder.toByteArray();
							   // if the ip is in ipv6 format, the decoded string is
							   // 16 bytes long. fast-forward 12 bytes
							   if (bytes.length == 16) {
								   for (var i:int = 0; i < 12; i++) {
									   bytes.readUnsignedByte();
								   }
							   }
							   ip_ipv4 = bytes.readUnsignedByte() + "." + bytes.readUnsignedByte() + "." + bytes.readUnsignedByte() + "." + bytes.readUnsignedByte();

							   _view.dockerApiEndpoint.text = "DOCKER_HOST=tcp://" + ip_ipv4;
							   _view.dockerLog.label = "https://" + ip_ipv4 + ":" + AppConstants.VCH_LOG_PORT;
							   continue;
						   }

						   if (keyName == AppConstants.DOCKER_PERSONALITY_ARGS_PATH) {
							   // port 2376 is used for tls, and 2375 for no-tls
							   isUsingTls = keyVal.indexOf(AppConstants.VCH_ENDPOINT_PORT_TLS) > -1;
							   continue;
						   }
					   }

					   // since the order in which list items are processed is not much guaranteed,
					   // we set the port for Docker API endpoint at the end of the loop
					   if (_view.dockerApiEndpoint.text != AppConstants.PLACEHOLDER_VAL) {
					       if (isUsingTls) {
							   _view.dockerApiEndpoint.text = _view.dockerApiEndpoint.text + ":" +
								   AppConstants.VCH_ENDPOINT_PORT_TLS;
					       } else {
							   _view.dockerApiEndpoint.text = _view.dockerApiEndpoint.text + ":" +
								   AppConstants.VCH_ENDPOINT_PORT_NO_TLS;
					       }
					   }
				   }
			   } else {
				   _view.isVch = false;
			   }
		   }
	   }

	   private function clearData() : void {
		   if(_view != null) {
			   // clear the UI data
			   _view.isVch = false;
			   _view.dockerApiEndpoint.text = new String(AppConstants.PLACEHOLDER_VAL);
			   _view.dockerLog.label = new String(AppConstants.PLACEHOLDER_VAL);
		   }
	   }
	}
}
