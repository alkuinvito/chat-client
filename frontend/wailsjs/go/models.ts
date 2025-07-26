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
	
	export class Response___chat_client_internal_chat_ChatRoom_ {
	    code: number;
	    data: chat.ChatRoom[];
	
	    static createFrom(source: any = {}) {
	        return new Response___chat_client_internal_chat_ChatRoom_(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.code = source["code"];
	        this.data = this.convertValues(source["data"], chat.ChatRoom);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class Response_chat_client_internal_user_UserProfile_ {
	    code: number;
	    data: user.UserProfile;
	
	    static createFrom(source: any = {}) {
	        return new Response_chat_client_internal_user_UserProfile_(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.code = source["code"];
	        this.data = this.convertValues(source["data"], user.UserProfile);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class Response_string_ {
	    code: number;
	    data: string;
	
	    static createFrom(source: any = {}) {
	        return new Response_string_(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.code = source["code"];
	        this.data = source["data"];
	    }
	}

}

export namespace user {
	
	export class UserProfile {
	    id: string;
	    username: string;
	
	    static createFrom(source: any = {}) {
	        return new UserProfile(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.username = source["username"];
	    }
	}

}

