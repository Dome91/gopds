import {ComponentFixture, TestBed} from '@angular/core/testing';

import {CatalogEntryComponent} from './catalog-entry.component';
import {CatalogEntry} from "../../models/catalog";

describe('CatalogEntryComponent', () => {
  let component: CatalogEntryComponent;
  let fixture: ComponentFixture<CatalogEntryComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      declarations: [ CatalogEntryComponent ]
    })
    .compileComponents();
  });

  beforeEach(() => {
    fixture = TestBed.createComponent(CatalogEntryComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });

  it('should set correct download link', () => {
    component.catalogEntry = new CatalogEntry('id1', 'name1', false, 'CBZ', '');
    component.ngOnInit();
    expect(component.downloadURL).toBe('http://localhost:3000/api/v1/catalog/id1/download');
  });

  it('should set correct cover link', () => {
    component.catalogEntry = new CatalogEntry('id1', 'name1', false, 'CBZ', '');
    component.ngOnInit();
    expect(component.coverURL).toBe('assets/card.svg');

    component.catalogEntry.cover = "cover1"
    component.ngOnInit();
    expect(component.coverURL).toBe('http://localhost:3000/covers/cover1');
  });
});
