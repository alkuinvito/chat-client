export namespace chat {
	
	export class ChatRoom {
	    peer_name: string;
	    ip: string;
	
	    static createFrom(source: any = {}) {
	        return new ChatRoom(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.peer_name = source["peer_name"];
	        this.ip = source["ip"];
	    }
	}

}

