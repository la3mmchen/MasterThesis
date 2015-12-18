<?php

class Object {
  public $UniqId;
  public $Id;
  public $Name;

  public function __construct() {
      $this->UniqId = rand(100000, 999999);
  }

  public function setName($name) {
    $this->Name = $name;
    return 0;
  }

  public function setId($id) {
    $this->Id = $id;
    return 0;
  }

  public function getJson() {
    return json_encode($this);
  }

  public function get($id) {
    return $this;
  }
}

?>
