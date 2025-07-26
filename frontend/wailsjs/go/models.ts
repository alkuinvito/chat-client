export namespace auth {
	
	export class UserModel {
	    username: string;
	
	    static createFrom(source: any = {}) {
	        return new UserModel(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.username = source["username"];
	    }
	}

}

export namespace chat {
	
	export class ChatMessage {
	    sender: string;
	    message: string;
	
	    static createFrom(source: any = {}) {
	        return new ChatMessage(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.sender = source["sender"];
	        this.message = source["message"];
	    }
	}
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

export namespace response {
	
	export class Response {
	    code: number;
	    message: string;
	
	    static createFrom(source: any = {}) {
	        return new Response(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.code = source["code"];
	        this.message = source["message"];
	    }
	}

}

