import {TestBed} from '@angular/core/testing';

import {HttpErrorInterceptor} from './http-error.interceptor';
import {HttpErrorResponse, HttpEvent, HttpHandler, HttpRequest} from "@angular/common/http";
import {Observable, throwError} from "rxjs";
import {ToastrModule} from "ngx-toastr";
import createSpyObj = jasmine.createSpyObj;
import {Router} from "@angular/router";
import SpyObj = jasmine.SpyObj;

describe('HttpErrorInterceptor', () => {

  let router: SpyObj<Router>;
  let interceptor: HttpErrorInterceptor

  beforeEach(() => {
    router = createSpyObj<Router>(['navigateByUrl']);
    TestBed.configureTestingModule({
      imports: [ToastrModule.forRoot()],
      providers: [
        HttpErrorInterceptor,
        {provide: Router, useValue: router}
      ]
    })
  });

  beforeEach(() => interceptor = TestBed.inject(HttpErrorInterceptor));

  it('should be created', () => expect(interceptor).toBeTruthy());

  it('should navigate to root for status unauthorized', () => {
    // @ts-ignore
    interceptor.intercept(null, new MockHttpHandler(401)).subscribe(
      () => {
      }, _ => {
      });
    expect(router.navigateByUrl).toHaveBeenCalledWith('/');
  });

  it('should navigate to catalog for status forbidden', () => {
    // @ts-ignore
    interceptor.intercept(null, new MockHttpHandler(403)).subscribe(
      () => {
      },
      _ => {
      });
    expect(router.navigateByUrl).toHaveBeenCalledWith('/catalog');
  });
});

class MockHttpHandler extends HttpHandler {

  constructor(private status: number) {
    super();
  }

  handle(req: HttpRequest<any>): Observable<HttpEvent<any>> {
    const httpErrorResponse = new HttpErrorResponse({
      status: this.status,
    });

    return throwError(httpErrorResponse);
  }
}
