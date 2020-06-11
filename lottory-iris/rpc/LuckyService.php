<?php
namespace rpc;
/**
 * Autogenerated by Thrift Compiler (0.11.0)
 *
 * DO NOT EDIT UNLESS YOU ARE SURE THAT YOU KNOW WHAT YOU ARE DOING
 *  @generated
 */
use Thrift\Base\TBase;
use Thrift\Type\TType;
use Thrift\Type\TMessageType;
use Thrift\Exception\TException;
use Thrift\Exception\TProtocolException;
use Thrift\Protocol\TProtocol;
use Thrift\Protocol\TBinaryProtocolAccelerated;
use Thrift\Exception\TApplicationException;


interface LuckyServiceIf {
  /**
   * @param int $uid
   * @param string $username
   * @param string $ip
   * @param int $now
   * @param string $app
   * @param string $sign
   * @return \rpc\DataResult
   */
  public function DoLucky($uid, $username, $ip, $now, $app, $sign);
  /**
   * @param int $uid
   * @param string $username
   * @param string $ip
   * @param int $now
   * @param string $app
   * @param string $sign
   * @return \rpc\DataGiftPrize[]
   */
  public function MyPrizeList($uid, $username, $ip, $now, $app, $sign);
}


class LuckyServiceClient implements \rpc\LuckyServiceIf {
  protected $input_ = null;
  protected $output_ = null;

  protected $seqid_ = 0;

  public function __construct($input, $output=null) {
    $this->input_ = $input;
    $this->output_ = $output ? $output : $input;
  }

  public function DoLucky($uid, $username, $ip, $now, $app, $sign)
  {
    $this->send_DoLucky($uid, $username, $ip, $now, $app, $sign);
    return $this->recv_DoLucky();
  }

  public function send_DoLucky($uid, $username, $ip, $now, $app, $sign)
  {
    $args = new \rpc\LuckyService_DoLucky_args();
    $args->uid = $uid;
    $args->username = $username;
    $args->ip = $ip;
    $args->now = $now;
    $args->app = $app;
    $args->sign = $sign;
    $bin_accel = ($this->output_ instanceof TBinaryProtocolAccelerated) && function_exists('thrift_protocol_write_binary');
    if ($bin_accel)
    {
      thrift_protocol_write_binary($this->output_, 'DoLucky', TMessageType::CALL, $args, $this->seqid_, $this->output_->isStrictWrite());
    }
    else
    {
      $this->output_->writeMessageBegin('DoLucky', TMessageType::CALL, $this->seqid_);
      $args->write($this->output_);
      $this->output_->writeMessageEnd();
      $this->output_->getTransport()->flush();
    }
  }

  public function recv_DoLucky()
  {
    $bin_accel = ($this->input_ instanceof TBinaryProtocolAccelerated) && function_exists('thrift_protocol_read_binary');
    if ($bin_accel) $result = thrift_protocol_read_binary($this->input_, '\rpc\LuckyService_DoLucky_result', $this->input_->isStrictRead());
    else
    {
      $rseqid = 0;
      $fname = null;
      $mtype = 0;

      $this->input_->readMessageBegin($fname, $mtype, $rseqid);
      if ($mtype == TMessageType::EXCEPTION) {
        $x = new TApplicationException();
        $x->read($this->input_);
        $this->input_->readMessageEnd();
        throw $x;
      }
      $result = new \rpc\LuckyService_DoLucky_result();
      $result->read($this->input_);
      $this->input_->readMessageEnd();
    }
    if ($result->success !== null) {
      return $result->success;
    }
    throw new \Exception("DoLucky failed: unknown result");
  }

  public function MyPrizeList($uid, $username, $ip, $now, $app, $sign)
  {
    $this->send_MyPrizeList($uid, $username, $ip, $now, $app, $sign);
    return $this->recv_MyPrizeList();
  }

  public function send_MyPrizeList($uid, $username, $ip, $now, $app, $sign)
  {
    $args = new \rpc\LuckyService_MyPrizeList_args();
    $args->uid = $uid;
    $args->username = $username;
    $args->ip = $ip;
    $args->now = $now;
    $args->app = $app;
    $args->sign = $sign;
    $bin_accel = ($this->output_ instanceof TBinaryProtocolAccelerated) && function_exists('thrift_protocol_write_binary');
    if ($bin_accel)
    {
      thrift_protocol_write_binary($this->output_, 'MyPrizeList', TMessageType::CALL, $args, $this->seqid_, $this->output_->isStrictWrite());
    }
    else
    {
      $this->output_->writeMessageBegin('MyPrizeList', TMessageType::CALL, $this->seqid_);
      $args->write($this->output_);
      $this->output_->writeMessageEnd();
      $this->output_->getTransport()->flush();
    }
  }

  public function recv_MyPrizeList()
  {
    $bin_accel = ($this->input_ instanceof TBinaryProtocolAccelerated) && function_exists('thrift_protocol_read_binary');
    if ($bin_accel) $result = thrift_protocol_read_binary($this->input_, '\rpc\LuckyService_MyPrizeList_result', $this->input_->isStrictRead());
    else
    {
      $rseqid = 0;
      $fname = null;
      $mtype = 0;

      $this->input_->readMessageBegin($fname, $mtype, $rseqid);
      if ($mtype == TMessageType::EXCEPTION) {
        $x = new TApplicationException();
        $x->read($this->input_);
        $this->input_->readMessageEnd();
        throw $x;
      }
      $result = new \rpc\LuckyService_MyPrizeList_result();
      $result->read($this->input_);
      $this->input_->readMessageEnd();
    }
    if ($result->success !== null) {
      return $result->success;
    }
    throw new \Exception("MyPrizeList failed: unknown result");
  }

}


// HELPER FUNCTIONS AND STRUCTURES

class LuckyService_DoLucky_args {
  static $isValidate = false;

  static $_TSPEC = array(
    1 => array(
      'var' => 'uid',
      'isRequired' => false,
      'type' => TType::I64,
      ),
    2 => array(
      'var' => 'username',
      'isRequired' => false,
      'type' => TType::STRING,
      ),
    3 => array(
      'var' => 'ip',
      'isRequired' => false,
      'type' => TType::STRING,
      ),
    4 => array(
      'var' => 'now',
      'isRequired' => false,
      'type' => TType::I64,
      ),
    5 => array(
      'var' => 'app',
      'isRequired' => false,
      'type' => TType::STRING,
      ),
    6 => array(
      'var' => 'sign',
      'isRequired' => false,
      'type' => TType::STRING,
      ),
    );

  /**
   * @var int
   */
  public $uid = null;
  /**
   * @var string
   */
  public $username = null;
  /**
   * @var string
   */
  public $ip = null;
  /**
   * @var int
   */
  public $now = null;
  /**
   * @var string
   */
  public $app = null;
  /**
   * @var string
   */
  public $sign = null;

  public function __construct($vals=null) {
    if (is_array($vals)) {
      if (isset($vals['uid'])) {
        $this->uid = $vals['uid'];
      }
      if (isset($vals['username'])) {
        $this->username = $vals['username'];
      }
      if (isset($vals['ip'])) {
        $this->ip = $vals['ip'];
      }
      if (isset($vals['now'])) {
        $this->now = $vals['now'];
      }
      if (isset($vals['app'])) {
        $this->app = $vals['app'];
      }
      if (isset($vals['sign'])) {
        $this->sign = $vals['sign'];
      }
    }
  }

  public function getName() {
    return 'LuckyService_DoLucky_args';
  }

  public function read($input)
  {
    $xfer = 0;
    $fname = null;
    $ftype = 0;
    $fid = 0;
    $xfer += $input->readStructBegin($fname);
    while (true)
    {
      $xfer += $input->readFieldBegin($fname, $ftype, $fid);
      if ($ftype == TType::STOP) {
        break;
      }
      switch ($fid)
      {
        case 1:
          if ($ftype == TType::I64) {
            $xfer += $input->readI64($this->uid);
          } else {
            $xfer += $input->skip($ftype);
          }
          break;
        case 2:
          if ($ftype == TType::STRING) {
            $xfer += $input->readString($this->username);
          } else {
            $xfer += $input->skip($ftype);
          }
          break;
        case 3:
          if ($ftype == TType::STRING) {
            $xfer += $input->readString($this->ip);
          } else {
            $xfer += $input->skip($ftype);
          }
          break;
        case 4:
          if ($ftype == TType::I64) {
            $xfer += $input->readI64($this->now);
          } else {
            $xfer += $input->skip($ftype);
          }
          break;
        case 5:
          if ($ftype == TType::STRING) {
            $xfer += $input->readString($this->app);
          } else {
            $xfer += $input->skip($ftype);
          }
          break;
        case 6:
          if ($ftype == TType::STRING) {
            $xfer += $input->readString($this->sign);
          } else {
            $xfer += $input->skip($ftype);
          }
          break;
        default:
          $xfer += $input->skip($ftype);
          break;
      }
      $xfer += $input->readFieldEnd();
    }
    $xfer += $input->readStructEnd();
    return $xfer;
  }

  public function write($output) {
    $xfer = 0;
    $xfer += $output->writeStructBegin('LuckyService_DoLucky_args');
    if ($this->uid !== null) {
      $xfer += $output->writeFieldBegin('uid', TType::I64, 1);
      $xfer += $output->writeI64($this->uid);
      $xfer += $output->writeFieldEnd();
    }
    if ($this->username !== null) {
      $xfer += $output->writeFieldBegin('username', TType::STRING, 2);
      $xfer += $output->writeString($this->username);
      $xfer += $output->writeFieldEnd();
    }
    if ($this->ip !== null) {
      $xfer += $output->writeFieldBegin('ip', TType::STRING, 3);
      $xfer += $output->writeString($this->ip);
      $xfer += $output->writeFieldEnd();
    }
    if ($this->now !== null) {
      $xfer += $output->writeFieldBegin('now', TType::I64, 4);
      $xfer += $output->writeI64($this->now);
      $xfer += $output->writeFieldEnd();
    }
    if ($this->app !== null) {
      $xfer += $output->writeFieldBegin('app', TType::STRING, 5);
      $xfer += $output->writeString($this->app);
      $xfer += $output->writeFieldEnd();
    }
    if ($this->sign !== null) {
      $xfer += $output->writeFieldBegin('sign', TType::STRING, 6);
      $xfer += $output->writeString($this->sign);
      $xfer += $output->writeFieldEnd();
    }
    $xfer += $output->writeFieldStop();
    $xfer += $output->writeStructEnd();
    return $xfer;
  }

}

class LuckyService_DoLucky_result {
  static $isValidate = false;

  static $_TSPEC = array(
    0 => array(
      'var' => 'success',
      'isRequired' => false,
      'type' => TType::STRUCT,
      'class' => '\rpc\DataResult',
      ),
    );

  /**
   * @var \rpc\DataResult
   */
  public $success = null;

  public function __construct($vals=null) {
    if (is_array($vals)) {
      if (isset($vals['success'])) {
        $this->success = $vals['success'];
      }
    }
  }

  public function getName() {
    return 'LuckyService_DoLucky_result';
  }

  public function read($input)
  {
    $xfer = 0;
    $fname = null;
    $ftype = 0;
    $fid = 0;
    $xfer += $input->readStructBegin($fname);
    while (true)
    {
      $xfer += $input->readFieldBegin($fname, $ftype, $fid);
      if ($ftype == TType::STOP) {
        break;
      }
      switch ($fid)
      {
        case 0:
          if ($ftype == TType::STRUCT) {
            $this->success = new \rpc\DataResult();
            $xfer += $this->success->read($input);
          } else {
            $xfer += $input->skip($ftype);
          }
          break;
        default:
          $xfer += $input->skip($ftype);
          break;
      }
      $xfer += $input->readFieldEnd();
    }
    $xfer += $input->readStructEnd();
    return $xfer;
  }

  public function write($output) {
    $xfer = 0;
    $xfer += $output->writeStructBegin('LuckyService_DoLucky_result');
    if ($this->success !== null) {
      if (!is_object($this->success)) {
        throw new TProtocolException('Bad type in structure.', TProtocolException::INVALID_DATA);
      }
      $xfer += $output->writeFieldBegin('success', TType::STRUCT, 0);
      $xfer += $this->success->write($output);
      $xfer += $output->writeFieldEnd();
    }
    $xfer += $output->writeFieldStop();
    $xfer += $output->writeStructEnd();
    return $xfer;
  }

}

class LuckyService_MyPrizeList_args {
  static $isValidate = false;

  static $_TSPEC = array(
    1 => array(
      'var' => 'uid',
      'isRequired' => false,
      'type' => TType::I64,
      ),
    2 => array(
      'var' => 'username',
      'isRequired' => false,
      'type' => TType::STRING,
      ),
    3 => array(
      'var' => 'ip',
      'isRequired' => false,
      'type' => TType::STRING,
      ),
    4 => array(
      'var' => 'now',
      'isRequired' => false,
      'type' => TType::I64,
      ),
    5 => array(
      'var' => 'app',
      'isRequired' => false,
      'type' => TType::STRING,
      ),
    6 => array(
      'var' => 'sign',
      'isRequired' => false,
      'type' => TType::STRING,
      ),
    );

  /**
   * @var int
   */
  public $uid = null;
  /**
   * @var string
   */
  public $username = null;
  /**
   * @var string
   */
  public $ip = null;
  /**
   * @var int
   */
  public $now = null;
  /**
   * @var string
   */
  public $app = null;
  /**
   * @var string
   */
  public $sign = null;

  public function __construct($vals=null) {
    if (is_array($vals)) {
      if (isset($vals['uid'])) {
        $this->uid = $vals['uid'];
      }
      if (isset($vals['username'])) {
        $this->username = $vals['username'];
      }
      if (isset($vals['ip'])) {
        $this->ip = $vals['ip'];
      }
      if (isset($vals['now'])) {
        $this->now = $vals['now'];
      }
      if (isset($vals['app'])) {
        $this->app = $vals['app'];
      }
      if (isset($vals['sign'])) {
        $this->sign = $vals['sign'];
      }
    }
  }

  public function getName() {
    return 'LuckyService_MyPrizeList_args';
  }

  public function read($input)
  {
    $xfer = 0;
    $fname = null;
    $ftype = 0;
    $fid = 0;
    $xfer += $input->readStructBegin($fname);
    while (true)
    {
      $xfer += $input->readFieldBegin($fname, $ftype, $fid);
      if ($ftype == TType::STOP) {
        break;
      }
      switch ($fid)
      {
        case 1:
          if ($ftype == TType::I64) {
            $xfer += $input->readI64($this->uid);
          } else {
            $xfer += $input->skip($ftype);
          }
          break;
        case 2:
          if ($ftype == TType::STRING) {
            $xfer += $input->readString($this->username);
          } else {
            $xfer += $input->skip($ftype);
          }
          break;
        case 3:
          if ($ftype == TType::STRING) {
            $xfer += $input->readString($this->ip);
          } else {
            $xfer += $input->skip($ftype);
          }
          break;
        case 4:
          if ($ftype == TType::I64) {
            $xfer += $input->readI64($this->now);
          } else {
            $xfer += $input->skip($ftype);
          }
          break;
        case 5:
          if ($ftype == TType::STRING) {
            $xfer += $input->readString($this->app);
          } else {
            $xfer += $input->skip($ftype);
          }
          break;
        case 6:
          if ($ftype == TType::STRING) {
            $xfer += $input->readString($this->sign);
          } else {
            $xfer += $input->skip($ftype);
          }
          break;
        default:
          $xfer += $input->skip($ftype);
          break;
      }
      $xfer += $input->readFieldEnd();
    }
    $xfer += $input->readStructEnd();
    return $xfer;
  }

  public function write($output) {
    $xfer = 0;
    $xfer += $output->writeStructBegin('LuckyService_MyPrizeList_args');
    if ($this->uid !== null) {
      $xfer += $output->writeFieldBegin('uid', TType::I64, 1);
      $xfer += $output->writeI64($this->uid);
      $xfer += $output->writeFieldEnd();
    }
    if ($this->username !== null) {
      $xfer += $output->writeFieldBegin('username', TType::STRING, 2);
      $xfer += $output->writeString($this->username);
      $xfer += $output->writeFieldEnd();
    }
    if ($this->ip !== null) {
      $xfer += $output->writeFieldBegin('ip', TType::STRING, 3);
      $xfer += $output->writeString($this->ip);
      $xfer += $output->writeFieldEnd();
    }
    if ($this->now !== null) {
      $xfer += $output->writeFieldBegin('now', TType::I64, 4);
      $xfer += $output->writeI64($this->now);
      $xfer += $output->writeFieldEnd();
    }
    if ($this->app !== null) {
      $xfer += $output->writeFieldBegin('app', TType::STRING, 5);
      $xfer += $output->writeString($this->app);
      $xfer += $output->writeFieldEnd();
    }
    if ($this->sign !== null) {
      $xfer += $output->writeFieldBegin('sign', TType::STRING, 6);
      $xfer += $output->writeString($this->sign);
      $xfer += $output->writeFieldEnd();
    }
    $xfer += $output->writeFieldStop();
    $xfer += $output->writeStructEnd();
    return $xfer;
  }

}

class LuckyService_MyPrizeList_result {
  static $isValidate = false;

  static $_TSPEC = array(
    0 => array(
      'var' => 'success',
      'isRequired' => false,
      'type' => TType::LST,
      'etype' => TType::STRUCT,
      'elem' => array(
        'type' => TType::STRUCT,
        'class' => '\rpc\DataGiftPrize',
        ),
      ),
    );

  /**
   * @var \rpc\DataGiftPrize[]
   */
  public $success = null;

  public function __construct($vals=null) {
    if (is_array($vals)) {
      if (isset($vals['success'])) {
        $this->success = $vals['success'];
      }
    }
  }

  public function getName() {
    return 'LuckyService_MyPrizeList_result';
  }

  public function read($input)
  {
    $xfer = 0;
    $fname = null;
    $ftype = 0;
    $fid = 0;
    $xfer += $input->readStructBegin($fname);
    while (true)
    {
      $xfer += $input->readFieldBegin($fname, $ftype, $fid);
      if ($ftype == TType::STOP) {
        break;
      }
      switch ($fid)
      {
        case 0:
          if ($ftype == TType::LST) {
            $this->success = array();
            $_size0 = 0;
            $_etype3 = 0;
            $xfer += $input->readListBegin($_etype3, $_size0);
            for ($_i4 = 0; $_i4 < $_size0; ++$_i4)
            {
              $elem5 = null;
              $elem5 = new \rpc\DataGiftPrize();
              $xfer += $elem5->read($input);
              $this->success []= $elem5;
            }
            $xfer += $input->readListEnd();
          } else {
            $xfer += $input->skip($ftype);
          }
          break;
        default:
          $xfer += $input->skip($ftype);
          break;
      }
      $xfer += $input->readFieldEnd();
    }
    $xfer += $input->readStructEnd();
    return $xfer;
  }

  public function write($output) {
    $xfer = 0;
    $xfer += $output->writeStructBegin('LuckyService_MyPrizeList_result');
    if ($this->success !== null) {
      if (!is_array($this->success)) {
        throw new TProtocolException('Bad type in structure.', TProtocolException::INVALID_DATA);
      }
      $xfer += $output->writeFieldBegin('success', TType::LST, 0);
      {
        $output->writeListBegin(TType::STRUCT, count($this->success));
        {
          foreach ($this->success as $iter6)
          {
            $xfer += $iter6->write($output);
          }
        }
        $output->writeListEnd();
      }
      $xfer += $output->writeFieldEnd();
    }
    $xfer += $output->writeFieldStop();
    $xfer += $output->writeStructEnd();
    return $xfer;
  }

}


