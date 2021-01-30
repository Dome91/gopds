import {Injectable} from '@angular/core';
import {HttpClient} from "@angular/common/http";
import {Observable} from "rxjs";
import {CatalogEntriesInPage, CatalogEntry} from "../models/catalog";

@Injectable({
  providedIn: 'root'
})
export class CatalogService {

  constructor(private http: HttpClient) {
  }

  fetchInPage(page: number, pageSize: number, id: string): Observable<CatalogEntriesInPage> {
    return this.http.get<GetCatalogEntriesInPage>('/api/v1/catalog', {
      params: {
        page: page.toString(),
        pageSize: pageSize.toString(),
        id
      }
    })
  }
}

interface GetCatalogEntriesInPage {
  total: number;
  catalogEntries: CatalogEntry[];
}
