import {TestBed} from '@angular/core/testing';

import {FetchAllSourcesResponse, FetchSourceResponse, SourceService} from './source.service';
import {HttpClientTestingModule, HttpTestingController} from "@angular/common/http/testing";
import {Source} from "../models/source";

describe('SourceService', () => {
  let service: SourceService;
  let http: HttpTestingController

  beforeEach(() => {
    TestBed.configureTestingModule({
      imports: [HttpClientTestingModule]
    });
    service = TestBed.inject(SourceService);
    http = TestBed.inject(HttpTestingController)
  });

  it('should be created', () => {
    expect(service).toBeTruthy();
  });

  it('should create source', () => {
    const source = new Source("", 'source1', 'path/source1');
    service.create(source).subscribe();
    const req = http.expectOne('/api/v1/sources');
    expect(req.request.method).toBe('POST');
    req.flush(null);
  });

  it('should fetch all sources', () => {
    const source1 = new Source('id1', 'name1', 'path1')
    const source2 = new Source('id2', 'name2', 'path2')

    service.fetchAll().subscribe(
      response => {
        expect(response).toContain(source1);
        expect(response).toContain(source2);
      });

    const req = http.expectOne('/api/v1/sources');
    expect(req.request.method).toBe('GET');
    req.flush(new FetchAllSourcesResponse([
      new FetchSourceResponse(source1.id, source1.name, source1.path),
      new FetchSourceResponse(source2.id, source2.name, source2.path)]
    ));
  });

  it('should delete a source', () => {
    service.delete('id1').subscribe(() => {
    });

    const req = http.expectOne('/api/v1/sources/id1');
    expect(req.request.method).toBe('DELETE');
    req.flush(null);
  });

  it('should sync a source', () => {
    service.synchronize('id1').subscribe(() => {
    });

    const req = http.expectOne('/api/v1/sources/id1/sync');
    expect(req.request.method).toBe('PUT');
    req.flush(null);
  });
});
