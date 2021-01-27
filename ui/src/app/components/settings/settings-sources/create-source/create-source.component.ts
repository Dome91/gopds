import {Component, OnInit} from '@angular/core';
import {Source} from "../../../../models/source";
import {SourceService} from "../../../../services/source.service";
import {Router} from "@angular/router";

@Component({
  selector: 'app-create-source',
  templateUrl: './create-source.component.html',
  styleUrls: ['./create-source.component.sass']
})
export class CreateSourceComponent implements OnInit {

  source: Source;

  constructor(private sourceService: SourceService,
              private router: Router) {
    this.source = new Source('', '', '');
  }

  ngOnInit(): void {
  }

  create() {
    this.sourceService.create(this.source).subscribe(
      () => this.router.navigateByUrl('/settings'),
    );
  }
}
