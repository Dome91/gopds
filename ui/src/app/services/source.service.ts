import {Injectable} from '@angular/core';
import {HttpClient} from "@angular/common/http";
import {Source} from "../models/source";
import {Observable} from "rxjs";
import {map} from "rxjs/operators";

@Injectable({
  providedIn: 'root'
})
export class SourceService {

  constructor(private http: HttpClient) {
  }

  create(source: Source): Observable<null> {
    return this.http.post<null>('/api/v1/sources', {
      'name': source.name,
      'path': source.path
    });
  }

  fetchAll(): Observable<Source[]> {
    return this.http.get<FetchAllSourcesResponse>('/api/v1/sources')
      .pipe(
        map((response: FetchAllSourcesResponse) => response.sources.map(value => new Source(value.id, value.name, value.path)))
      );
  }

  delete(id: string): Observable<null> {
    return this.http.delete<null>(`/api/v1/sources/${id}`);
  }

  synchronize(id: string): Observable<null> {
    return this.http.put<null>(`/api/v1/sources/${id}/sync`, null);
  }
}

export class FetchAllSourcesResponse {
  constructor(public sources: FetchSourceResponse[]) {
  }
}

export class FetchSourceResponse {
  constructor(public id: string,
              public name: string,
              public path: string) {
  }
}
