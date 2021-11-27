import osmtogeojson from "osmtogeojson";
import { point } from "@turf/helpers";
import nearestPointOnLine from "@turf/nearest-point-on-line";
import toInteger from "lodash/toInteger";
import split from "lodash/split";
import has from "lodash/has";

export const GET_DATA = "GET_DATA";
export const SET_START_TIME = "SET_START_TIME";
export const SET_END_TIME = "SET_END_TIME";
export const SET_TOKEN = "SET_TOKEN";

const getLatLngBounds = dataJSON => {
  let result = {
    minLat: Number.MAX_VALUE,
    minLng: Number.MAX_VALUE,
    maxLat: -Number.MAX_VALUE,
    maxLng: -Number.MAX_VALUE
  };

  dataJSON.forEach(d => {
    if (d.Lat < result.minLat) {
      result.minLat = d.Lat;
    }
    if (d.Lat > result.maxLat) {
      result.maxLat = d.Lat;
    }

    if (d.Lng < result.minLng) {
      result.minLng = d.Lng;
    }
    if (d.Lng > result.maxLng) {
      result.maxLng = d.Lng;
    }
  });

  return result;
};

const handleHTTPResponseError = res => {
  if (res.status >= 400) {
    return res.text().then(text => {
      throw new Error(`API request error: ${res.status}: ${text}`);
    });
  }
  return res;
};

const MISSING = -999;

export const getData = (token, startTime, endTime) => {
  let dataJSON;
  return {
    type: GET_DATA,
    payload: fetch(
      `/api/records/${token}?startTime=${startTime.format()}&endTime=${endTime.format()}`
    )
      .then(handleHTTPResponseError)
      .then(res => res.json())
      .then(json => {
        dataJSON = json;
        const filteredDataJSON = dataJSON.filter(
          d => d.Lat !== MISSING && d.Lng !== MISSING
        );

        let minLat, minLng, maxLat, maxLng;
        if (filteredDataJSON.length > 1) {
          ({ minLat, minLng, maxLat, maxLng } = getLatLngBounds(
            filteredDataJSON
          ));
        } else {
          [minLat, minLng, maxLat, maxLng] = [0, 0, 0, 0];
        }
        return fetch(`http://overpass-api.de/api/interpreter`, {
          headers: {
            "content-type": "application/x-www-form-urlencoded"
          },
          method: "POST",
          body: `[out:json];way[highway][maxspeed](${minLat}, ${minLng}, ${maxLat}, ${maxLng});out body;>;out skel qt;`
        });
      })
      .then(handleHTTPResponseError)
      .then(res => res.json())
      .then(resJSON => {
        const roadGeoJSON = osmtogeojson(resJSON);
        dataJSON.forEach(d => {
          d.SpeedLimit = MISSING;
          if (d.Lng === MISSING || d.Lat === MISSING) {
            return;
          }
          const dPoint = point([d.Lng, d.Lat]);
          let closestDistance = Number.MAX_VALUE;
          let closestSpeedLimit;
          roadGeoJSON.features.forEach(f => {
            const lineString = f.geometry;
            const snapped = nearestPointOnLine(lineString, dPoint);
            if (
              snapped.properties.dist < closestDistance &&
              has(f.properties, "maxspeed")
            ) {
              closestDistance = snapped.properties.dist;
              closestSpeedLimit = toInteger(
                split(f.properties.maxspeed, " ")[0]
              );
            }
          });
          d.SpeedLimit = closestSpeedLimit;
        });
        return dataJSON;
      })
  };
};

export const setStartTime = startTime => ({
  type: SET_START_TIME,
  payload: startTime
});

export const setEndTime = endTime => ({
  type: SET_END_TIME,
  payload: endTime
});

export const setToken = token => ({
  type: SET_TOKEN,
  payload: token
});
