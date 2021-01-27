import { TestBed } from '@angular/core/testing';

import { LoginService } from './login.service';
import {HttpClientTestingModule, HttpTestingController} from "@angular/common/http/testing";

describe('LoginService', () => {
  let service: LoginService;
  let client: HttpTestingController

  beforeEach(() => {
    TestBed.configureTestingModule({
      imports: [HttpClientTestingModule]
    });

    service = TestBed.inject(LoginService);
    client = TestBed.inject(HttpTestingController)
  });

  it('should be created', () => {
    expect(service).toBeTruthy();
  });

  it('should log user in', () => {
    service.login('username', 'password').subscribe();
    const req = client.expectOne('/api/v1/login');
    expect(req.request.method).toBe('POST');
    req.flush(null);
  });

  it('should check if user is logged in', () => {
    service.isLoggedIn().subscribe();
    const req = client.expectOne('/api/v1/login');
    expect(req.request.method).toBe('GET');
    req.flush(null);
  })

});
