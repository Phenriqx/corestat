export namespace cpu {
	
	export class InfoStat {
	    cpu: number;
	    vendorId: string;
	    family: string;
	    model: string;
	    stepping: number;
	    physicalId: string;
	    coreId: string;
	    cores: number;
	    modelName: string;
	    mhz: number;
	    cacheSize: number;
	    flags: string[];
	    microcode: string;
	
	    static createFrom(source: any = {}) {
	        return new InfoStat(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.cpu = source["cpu"];
	        this.vendorId = source["vendorId"];
	        this.family = source["family"];
	        this.model = source["model"];
	        this.stepping = source["stepping"];
	        this.physicalId = source["physicalId"];
	        this.coreId = source["coreId"];
	        this.cores = source["cores"];
	        this.modelName = source["modelName"];
	        this.mhz = source["mhz"];
	        this.cacheSize = source["cacheSize"];
	        this.flags = source["flags"];
	        this.microcode = source["microcode"];
	    }
	}

}

export namespace main {
	
	export class CpuInformation {
	    CPUPercent: number[];
	    CPUCores: number;
	    CPUInfo: cpu.InfoStat[];
	
	    static createFrom(source: any = {}) {
	        return new CpuInformation(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.CPUPercent = source["CPUPercent"];
	        this.CPUCores = source["CPUCores"];
	        this.CPUInfo = this.convertValues(source["CPUInfo"], cpu.InfoStat);
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

}

