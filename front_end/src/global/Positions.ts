export type Position = { latitude: number; longitude: number };

export class PositionClass {
  getPosition(): Promise<Position> {
    return new Promise<Position>((resolve) => {
      if (!navigator.geolocation) {
        throw "Geolocation is not supported by your browser";
      }

      navigator.geolocation.getCurrentPosition(
        (position: GeolocationPosition): void => {
          const res: Position = {
            latitude: position.coords.latitude,
            longitude: position.coords.longitude,
          };
          // console.log("res before : ", res);
          resolve(res);
        },
        () => {
          throw "Unable to retrieve your location";
        },
        {
          enableHighAccuracy: true,
          // maximumAge?: number;
          // timeout?: number;
        }
      );
    });
  }
}
