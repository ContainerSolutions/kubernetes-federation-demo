/**
 * Simulate backend interactions during development
 */
(function() {

    Backend.adminUrl = 'http://localhost:8080';
    
    let promise = { 
        done: function(fun) { window.setTimeout(fun, randInt(50) + 50); return promise; },
        fail: function(fun) { return promise; }
    };

    let recentlyChanged = [];

    Backend.enable = function(name) {
        cannedDatacenters[name].Ready = true;
        recentlyChanged.push(name);
        return promise;
    };

    Backend.disable = function(name) {
        cannedDatacenters[name].Ready = false;
        recentlyChanged.push(name);
        return promise;
    };

    Backend.startTraffic = function(sourceName) {
        trafficSources[sourceName] = true;
        return promise;
    };

    Backend.stopTraffic = function(sourceName) {
        trafficSources[sourceName] = false;
        return promise;
    };

    Backend.isActiveTrafficSource = function(name) {
        return !!trafficSources[name];
    };

    let ticker = null, traffic = {};

    Backend.init = function() {
        console.log("Initializing simulated backend");
        window.setTimeout(function() {
            Backend.updateDCs(Object.values(cannedDatacenters));
        }, 250);
        ticker = window.setInterval(function() {
            try {
                let sources = activeSources(), src, updates = [], updatedNames = [];
                while (src = sources.shift()) {
                    let dc = nearestHealthyDC(src);
                    if (!dc) continue;
                    
                    let update = cloneDC(dc);
                    update["Traffic"] = randomTraffic(src, dc);
                    updates.push(update);
                    updatedNames.push(update.IpAddress.dc);
                }
                let name;
                while(name = recentlyChanged.shift()) {
                    if (updatedNames.indexOf(name) == -1) {
                        updates.push(cloneDC(cannedDatacenters[name]));
                    }
                }
                if (updates.length) {
                    Backend.updateDCs(updates);
                }
            } catch(e) {
                Backend.stop();
                throw(e);
            }
        }, randInt(100) + 500);
    };

    Backend.stop = function() {
        console.log("Stopping simulated backend");
        window.clearInterval(ticker);
    };

    function cloneDC(dc) {
        return $.extend(true, {}, dc);
    }

    function randomTraffic(src, targetDc) {
        let key = src + ">" + targetDc.IpAddress.dc, ret = {};
        if (!traffic[key]) traffic[key] = 0;
        traffic[key] += (randInt(10) + 1);
        ret[src] = traffic[key];
        return ret;
    }

    function activeSources() {
        return Object.keys(trafficSources).filter(function(s) { return trafficSources[s] });
    }

    function nearestHealthyDC(src) {
        let healthy = Object.values(cannedDatacenters).filter(function(dc) { return dc.Ready; }).map(function(dc) { return dc.IpAddress.dc; });
        for (let d of nearestDCs[src]) if (healthy.indexOf(d) > -1) return cannedDatacenters[d];
    }

    function randInt(max) {
        return Math.floor((Math.random() * max))
    }

    let trafficSources = { "asia-northeast1": true, "us-central1": true, "us-west1": true };

    let nearestDCs = {
        "asia-northeast1": ["asia-northeast1", "asia-east1", "us-west1", "us-central1", "europe-west1", "us-east1"],
        "us-central1": ["us-central1", "us-west1", "us-east1", "europe-west1", "asia-northeast1", "asia-east1"],
        "us-west1": ["us-west1", "us-central1", "us-east1", "asia-northeast1", "asia-east1", "europe-west1"],
        "europe-west1": ["europe-west1", "us-east1", "us-central1", "asia-east1", "us-west1", "asia-northeast1"],
        "asia-east1": ["asia-east1", "asia-northeast1", "us-west1", "us-central1", "europe-west1", "us-east1"],
        "us-east1": ["us-east1", "us-central1", "us-west1", "europe-west1", "asia-northeast1", "asia-east1"]
    };

    const dcEurope = {
        "CloudProvider":2,"Name":"projects/185454036493/zones/europe-west1-c",
        "IpAddress":{"dc": "europe-west1", "ip":"35.187.10.78","loc":"50.470976, 3.864521","city":"Saint-Ghislain","region":"Ghlin","country":"Belgium"},
        "ClientIpAddress":null,"Ready":true,"Timestamp":"2017-01-26T13:47:53.730552094Z"
    };

    const dcAsia = {
        "CloudProvider":2,"Name":"projects/185454036493/zones/asia-east1-c",
        "IpAddress":{"dc": "asia-east1", "ip":"35.187.10.78","loc":"23.925895, 120.441405","city":"Changhua County","region":"","country":"Taiwan"},
        "ClientIpAddress":null,"Ready":true,"Timestamp":"2017-01-26T13:47:53.730552094Z"
    };

    const dcUS = {
        "CloudProvider":2,"Name":"projects/185454036493/zones/us-east1-c",
        "IpAddress":{"dc": "us-east1", "ip":"35.187.10.78","loc":"33.072657, -80.038877","city":"Berkeley County","region":"South Carolina","country":"USA"},
        "ClientIpAddress":null,"Ready":false,"Timestamp":"2017-01-26T13:47:53.730552094Z"
    };

    const cannedDatacenters = { "europe-west1": dcEurope, "asia-east1": dcAsia, "us-east1": dcUS };

})();
