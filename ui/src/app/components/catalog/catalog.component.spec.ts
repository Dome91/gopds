import {ComponentFixture, TestBed} from '@angular/core/testing';

import {CatalogComponent} from './catalog.component';
import {CatalogService} from "../../services/catalog.service";
import {of} from "rxjs";
import {MockNavbarComponent} from "../mocks/mock-navbar.components";
import {ToastrModule} from "ngx-toastr";
import {FontAwesomeModule} from "@fortawesome/angular-fontawesome";
import {InfiniteScrollModule} from "ngx-infinite-scroll";
import {ActivatedRoute} from "@angular/router";
import {CatalogEntriesInPage, CatalogEntry} from "../../models/catalog";
import SpyObj = jasmine.SpyObj;
import createSpyObj = jasmine.createSpyObj;

describe('CatalogComponent', () => {
  let component: CatalogComponent;
  let fixture: ComponentFixture<CatalogComponent>;
  let catalogService: SpyObj<CatalogService>;

  let catalogEntriesInPage = new CatalogEntriesInPage(48, [new CatalogEntry('id2', 'name2', true), new CatalogEntry('id3', 'name3', false)]);

  beforeEach(async () => {
    catalogService = createSpyObj<CatalogService>(['fetchInPage']);
    let route = {
      queryParams: of({id: 'id1'})
    };

    await TestBed.configureTestingModule({
      declarations: [CatalogComponent, MockNavbarComponent],
      imports: [ToastrModule.forRoot(), FontAwesomeModule, InfiniteScrollModule],
      providers: [
        {provide: CatalogService, useValue: catalogService},
        {provide: ActivatedRoute, useValue: route}
      ]
    })
      .compileComponents();
  });

  beforeEach(() => {
    catalogService.fetchInPage.and.returnValue(of(catalogEntriesInPage));
    fixture = TestBed.createComponent(CatalogComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });

  it('should fetch another page on scrolling', () => {
    component.onScroll();
    expect(component.page).toEqual(1);
    expect(catalogService.fetchInPage).toHaveBeenCalledWith(1, 24, 'id1');
    expect(component.catalogEntries.length).toEqual(4)
  });
});
