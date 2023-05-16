export namespace main {
	
	export class Options {
	    testHost: string;
	    testDNSHost: string;
	
	    static createFrom(source: any = {}) {
	        return new Options(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.testHost = source["testHost"];
	        this.testDNSHost = source["testDNSHost"];
	    }
	}

}

