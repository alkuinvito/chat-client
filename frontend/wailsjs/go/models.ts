export namespace chat {
	
	export class ChatMessage {
	    id: number;
	    sender: string;
	    message: string;
	    created_at: string;
	
	    static createFrom(source: any = {}) {
	        return new ChatMessage(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.sender = source["sender"];
	        this.message = source["message"];
	        this.created_at = source["created_at"];
	    }
	}
	export class SendMessageSchema {
	    sender: string;
	    message: string;
	
	    static createFrom(source: any = {}) {
	        return new SendMessageSchema(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.sender = source["sender"];
	        this.message = source["message"];
	    }
	}

}

export namespace discovery {
	
	export class PeerModel {
	    id: string;
	    username: string;
	    ip: string;
	
	    static createFrom(source: any = {}) {
	        return new PeerModel(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.username = source["username"];
	        this.ip = source["ip"];
	    }
	}

}

export namespace response {
	
	export class Response___chat_client_internal_chat_ChatMessage_ {
	    code: number;
	    data: chat.ChatMessage[];
	
	    static createFrom(source: any = {}) {
	        return new Response___chat_client_internal_chat_ChatMessage_(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.code = source["code"];
	        this.data = this.convertValues(source["data"], chat.ChatMessage);
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
	export class Response___chat_client_internal_discovery_PeerModel_ {
	    code: number;
	    data: discovery.PeerModel[];
	
	    static createFrom(source: any = {}) {
	        return new Response___chat_client_internal_discovery_PeerModel_(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.code = source["code"];
	        this.data = this.convertValues(source["data"], discovery.PeerModel);
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
	export class Response___chat_client_internal_user_ContactModel_ {
	    code: number;
	    data: user.ContactModel[];
	
	    static createFrom(source: any = {}) {
	        return new Response___chat_client_internal_user_ContactModel_(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.code = source["code"];
	        this.data = this.convertValues(source["data"], user.ContactModel);
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
	export class Response_chat_client_internal_chat_ChatMessage_ {
	    code: number;
	    data: chat.ChatMessage;
	
	    static createFrom(source: any = {}) {
	        return new Response_chat_client_internal_chat_ChatMessage_(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.code = source["code"];
	        this.data = this.convertValues(source["data"], chat.ChatMessage);
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
	
	export class ContactModel {
	    id: string;
	    username: string;
	    SharedKey: number[];
	
	    static createFrom(source: any = {}) {
	        return new ContactModel(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.username = source["username"];
	        this.SharedKey = source["SharedKey"];
	    }
	}
	export class InitPairSchema {
	    id: string;
	    username: string;
	    code: string;
	    pubkey: string;
	
	    static createFrom(source: any = {}) {
	        return new InitPairSchema(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.username = source["username"];
	        this.code = source["code"];
	        this.pubkey = source["pubkey"];
	    }
	}
	export class RequestPairSchema {
	    id: string;
	    username: string;
	    code: string;
	
	    static createFrom(source: any = {}) {
	        return new RequestPairSchema(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.username = source["username"];
	        this.code = source["code"];
	    }
	}
	export class ResponsePairSchema {
	    pubkey: string;
	
	    static createFrom(source: any = {}) {
	        return new ResponsePairSchema(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.pubkey = source["pubkey"];
	    }
	}
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

