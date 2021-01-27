import {ComponentFixture, TestBed, waitForAsync} from '@angular/core/testing';

import { SettingsSourcesComponent } from './settings-sources.component';
import {SourceService} from "../../../services/source.service";
import SpyObj = jasmine.SpyObj;
import {Source} from "../../../models/source";
import {BrowserAnimationsModule} from "@angular/platform-browser/animations";
import {FontAwesomeModule} from "@fortawesome/angular-fontawesome";
import {ToastrModule, ToastrService} from "ngx-toastr";
import {of} from "rxjs";
import {By} from "@angular/platform-browser";

describe('SettingsSourcesComponent', () => {
  let component: SettingsSourcesComponent;
  let fixture: ComponentFixture<SettingsSourcesComponent>;

  let sourceService: SpyObj<SourceService>;
  const source1 = new Source('id1', 'name1', 'path1');
  const source2 = new Source('id2', 'name2', 'path2');

  beforeEach(async () => {
    sourceService = jasmine.createSpyObj<SourceService>(['fetchAll', 'delete']);

    await TestBed.configureTestingModule({
      declarations: [SettingsSourcesComponent],
      imports: [BrowserAnimationsModule, ToastrModule.forRoot(), FontAwesomeModule],
      providers: [ToastrService,
        {provide: SourceService, useValue: sourceService}]
    })
      .compileComponents();
  });

  beforeEach(() => {
    sourceService.fetchAll.and.returnValue(of([source1, source2]));

    fixture = TestBed.createComponent(SettingsSourcesComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });

  it('should render sources fetched', waitForAsync(() => {
    fixture.whenStable().then(() => {
      let source1Div = fixture.debugElement.query(By.css('.source-id1'));
      let source2Div = fixture.debugElement.query(By.css('.source-id2'));

      const title1 = source1Div.children[0].nativeElement.textContent
      const title2 = source2Div.children[0].nativeElement.textContent
      expect(title1).toEqual('name1');
      expect(title2).toEqual('name2');

      expect(component.isSyncing.get(source1.id)).toBeFalse();
      expect(component.isSyncing.get(source2.id)).toBeFalse();
    });
  }));

  it('should delete a source', waitForAsync(() => {
    sourceService.delete.withArgs('id1').and.returnValue(of(null));
    sourceService.fetchAll.and.returnValue(of([source2]));

    fixture.whenStable().then(() => {
      const deleteButton = fixture.debugElement.query(By.css('#delete-id1'));
      deleteButton.triggerEventHandler('click', null);

      expect(component.sources.length).toBe(1);
      expect(component.sources).toContain(source2);
      expect(component.isSyncing.size).toEqual(1);
      expect(component.isSyncing.get(source2.id)).toBeFalse();
    });
  }));
});
