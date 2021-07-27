export type PositionPoint = { latitude: number; longitude: number };

// export class PositionClass {
//   getPosition(): Promise<PositionPoint> {
//     return new Promise<PositionPoint>((resolve,reject) => {
//       if (!navigator.geolocation) {
//         throw "Geolocation is not supported by your browser";
//       }
//       navigator.geolocation.getCurrentPosition(
//         (position: GeolocationPosition): void => {
//           const res: PositionPoint = {
//             latitude: position.coords.latitude,
//             longitude: position.coords.longitude,
//           };
//           resolve(res);
//         },
//         () => {
//           reject("Unable to retrieve your location")
//           // throw "Unable to retrieve your location";
//         },
//         {
//           enableHighAccuracy: true,
//           // maximumAge?: number;
//           // timeout?: number;
//         }
//       );
//     });
//   }
// }

export function GetPosition(): Promise<PositionPoint> {
  return new Promise<PositionPoint>((resolve) => {
    // DEGUG_MODE
    const startPoint = {
      latitude: 25.040056717110396,
      longitude: 121.51187490970621,
    };
    resolve(startPoint);

    // if (!navigator.geolocation) {
    //   throw "Geolocation is not supported by your browser";
    // }
    // navigator.geolocation.getCurrentPosition(
    //   (position: GeolocationPosition): void => {
    //     const res: PositionPoint = {
    //       latitude: position.coords.latitude,
    //       longitude: position.coords.longitude,
    //     };
    //     resolve(res);
    //   },
    //   () => {
    //     reject("Unable to retrieve your location");
    //     // throw "Unable to retrieve your location";
    //   },
    //   {
    //     enableHighAccuracy: true,
    //     // maximumAge?: number;
    //     // timeout?: number;
    //   }
    // );
  });
}
