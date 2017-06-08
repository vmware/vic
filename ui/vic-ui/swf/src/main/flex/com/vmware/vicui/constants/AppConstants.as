package com.vmware.vicui.constants {

	public class AppConstants {

		public static const VM_CONTAINER_NAME_PATH:String = "common/name";
		public static const VM_CONTAINER_IMAGE_PATH:String = "guestinfo.vice./repo";
		public static const VM_CONTAINER_PORTMAPPING:String = "guestinfo.vice./networks|bridge/ports~";
		public static const VCH_NAME_PATH:String = "init/common/name";
		public static const VCH_CLIENT_IP_PATH:String = "guestinfo.vice..init.networks|client.assigned.IP";
		public static const DOCKER_PERSONALITY_ARGS_PATH:String = "guestinfo.vice./init/sessions|docker-personality/cmd/Args~";
		public static const VCH_ENDPOINT_PORT_TLS:String = "2376";
		public static const VCH_ENDPOINT_PORT_NO_TLS:String = "2375";
		public static const VCH_LOG_PORT:String = "2378";
		public static const PLACEHOLDER_VAL:String = "-";

	}

}
