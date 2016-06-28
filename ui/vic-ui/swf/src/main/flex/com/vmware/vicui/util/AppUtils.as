package com.vmware.vicui.util {
	
	public class AppUtils {
		
		public static function findIndexOfValue( data:Array, value:String ):int {
			var length:int = data.length;
			var xint:int = -1;
			for(var i:int=0; i<length; i++) {
				
				var optionValue:String = data[i].value;
				var key:String = data[i].key.value as String;
				
				if(key == value){
					xint = i;
					break;
				}
				
			}
			return xint;
		}
		
	}
}