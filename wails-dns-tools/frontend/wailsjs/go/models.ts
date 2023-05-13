export namespace main {
	
	export class Options {
	    pdfHeight: string;
	    pdfWidth: string;
	    pdfExportPath: string;
	
	    static createFrom(source: any = {}) {
	        return new Options(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.pdfHeight = source["pdfHeight"];
	        this.pdfWidth = source["pdfWidth"];
	        this.pdfExportPath = source["pdfExportPath"];
	    }
	}

}

