import {TestBed} from '@angular/core/testing';

import {CatalogService} from './catalog.service';
import {HttpClientTestingModule, HttpTestingController} from "@angular/common/http/testing";
import {CatalogEntriesInPage, CatalogEntry} from "../models/catalog";
import {CatalogEntryComponent} from "../components/catalog-entry/catalog-entry.component";

describe('CatalogService', () => {
  let service: CatalogService;
  let http: HttpTestingController;

  beforeEach(() => {
    TestBed.configureTestingModule({
      imports: [HttpClientTestingModule],
      declarations: [CatalogEntryComponent]
    });
    service = TestBed.inject(CatalogService);
    http = TestBed.inject(HttpTestingController);
  });

  it('should be created', () => {
    expect(service).toBeTruthy();
  });

  it('should fetch catalog entries in page', () => {
    const entry1 = new CatalogEntry('id1', 'name1', true, "", "");
    const entry2 = new CatalogEntry('id2', 'name2', false, "CBZ", "cover1");

    service.fetchInPage(0, 24, 'id1').subscribe(
      (response: CatalogEntriesInPage) => {
        expect(response.catalogEntries[0]).toEqual(entry1);
        expect(response.catalogEntries[1]).toEqual(entry2);
        expect(response.total).toEqual(2);
      });

    const req = http.expectOne('/api/v1/catalog?page=0&pageSize=24&id=id1');
    expect(req.request.method).toBe('GET');
    req.flush({catalogEntries: [entry1, entry2], total: 2});
  });

});
